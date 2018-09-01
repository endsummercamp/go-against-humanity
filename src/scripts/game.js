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
            style.height = percentage * 100 + "%";
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
            <div style={{position: "absolute", top: 0, left: 0, bottom: 0,right: 0}}>
                <div style={style} className="vote-bg"></div>
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
const socket = new WebSocket(`ws://${document.location.hostname}:8080/ws?match=${MATCH_ID}`);
socket.onopen = function() {
    console.log("Opened socket.");
};

function getCardText(data) {
    const {NewCard: { text: cardText }} = data
    return cardText
}

function getCardTotals(data) {
    const { Totals: _totals } = data
    return _totals
}

let answers = [];
let totals = [];
socket.onmessage = function (e) {
    console.log("Received", e.data);
    const data = JSON.parse(e.data);
    const { Name: eventName } = data 
    let cardText;
    switch (eventName) {
    case "join_successful":
        answers = [];
        ReactDOM.render(<BlackRow />, blackrowDiv);
        ReactDOM.render(<WhiteRow answers={[]}/>, whiterowDiv);
        break;
    case "new_black":
        ReactDOM.render(<BlackRow card={<Card text={data.NewCard.text} black />}/>, blackrowDiv);
        break;
    case "new_white":
        cardText = getCardText(data)
        answers.push({ text: data.NewCard.text, total: 0, ID: data.NewCard.Id });
        ReactDOM.render(<WhiteRow answers={answers}/>, whiterowDiv);
        break;
    case "totals":
        totals = data.Totals
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