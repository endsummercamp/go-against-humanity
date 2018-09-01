if (!window.WebSocket) {
    alert("Your browser does not support WebSockets!")
}

class Card extends React.Component {
    submitVote() {
        console.log(this.props);
        alert("POST /vote/" + this.props.id);
    }
    render() {
        // https://stackoverflow.com/a/6040258
        // https://stackoverflow.com/a/8541575
        const baseHeight = 20.25; // in rems
        let style = {
            position: "absolute",
            bottom: 0,
            left: 0,
            width: "100%",
            height: "0%",
            backgroundColor: "#a9f16c",
            zIndex: -1,
        };
        if (this.props.total) {
            const percentage = this.props.total / this.props.sum;
            // const percText = 100*percentage + "%";
            style.height = percentage * baseHeight + "rem";
        };
        console.log(this.props, style);
        return <div className={"card card-" + (this.props.black ? "black" : "white")} onClick={() => this.submitVote()}>
            <div className="card-top">
                <div className="card-content">
                    {this.props.text}
                </div>
            </div>
            <div className="card-middle">
                {this.props.total || ""}
            </div>
            <div className="card-bottom">
                Cards Against Humanity
            </div>
            <div style={{position: "relative", /* width: 0, */ height: 0}}>
                <div style={style} class="vote-bg"></div>
            </div>
        </div>
    }
}

class BlackRow extends React.PureComponent {
    render() {
        return <div className="flex" id="blackrow">
            {("card" in this.props) ?
                this.props.card :
                <i>In attesa di una nuova carta nera...</i>
            }
        </div>;
    }
}

class WhiteRow extends React.Component {
    render() {
        let sum = 0;
        if (this.props.totals) {
            sum = this.props.totals.reduce((a, b) => a + b.Votes, 0);
        }
        return <div className="flex" id="blackrow">
            {
                this.props.answers.map((answer, i) => <Card text={answer.text} id={answer.ID} total={answer.total} sum={sum} key={i} />)
            }
        </div>;
    }
}

const blackrowDiv = document.getElementById("react-blackrow");
ReactDOM.render(<BlackRow />, blackrowDiv);
const whiterowDiv = document.getElementById("react-whiterow");
ReactDOM.render(<WhiteRow answers={[]} />, whiterowDiv);

const whiteRow = document.getElementById("whiterow");
// TODO: cambiare URL e numero match
const socket = new WebSocket("ws://localhost:8080/ws?match=0");
socket.onopen = function() {
    console.log("Opened socket.");
};
let answers = [];
let totals = [];
socket.onmessage = function (e) {
    console.log("Received", e.data);
    const {Name: eventName, NewCard: {Text: cardText, ID: cardID}, Totals: _totals} = JSON.parse(e.data);
    switch (eventName) {
    case "new_game":
        answers = [];
        ReactDOM.render(<BlackRow />, blackrowDiv);
        ReactDOM.render(<WhiteRow answers={[]}/>, whiterowDiv);
        break;
    case "new_black":
        ReactDOM.render(<BlackRow card={<Card text={cardText} black />}/>, blackrowDiv);
        break;
    case "new_white":
        answers.push({text: cardText, ID: cardID, total: 0});
        ReactDOM.render(<WhiteRow answers={answers}/>, whiterowDiv);
        break;
    case "totals":
        totals = _totals;
        for (const total of totals) {
            answers.find(a => a.ID == total.ID).total = total.Votes;
        }
        ReactDOM.render(<WhiteRow answers={answers} totals={totals}/>, whiterowDiv);
        break;
    default:
        alert("Unknown event " + eventName);
    }
}
socket.onclose = function () {
    console.log("Socket closed.");
}