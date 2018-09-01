if (!window.WebSocket) {
    alert("Your browser does not support WebSockets!")
}

// TODO: uncomment
// const IS_PLAYER = !!window.IS_PLAYER;
const IS_PLAYER = 1;

class Card extends React.Component {
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
        return <div className={"card card-" + (this.props.black ? "black" : "white")} onClick={this.props.onClick}>
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

class AnswersRow extends React.Component {
    render() {
        /* Expects:
           * a prop "answers", containing an array of {text, ID};
           * a prop "totals", containing an array of {Votes}.
         */
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

class MyCardsRow extends React.Component {
    submitCard(id) {
        const req = new XMLHttpRequest();
        req.open("POST", `/pick_card?card_id=${id}&match_id=${0}`);
        req.send();
    }
    render() {
        /* Expects:
           * a prop "cards", containing an array of {text, ID};
         */
        return <>
            {
                this.props.cards.map((answer, i) => <Card text={answer.text} id={answer.ID} onClick={() => this.submitCard(answer.ID)} key={i} />)
            }
        </>;
    }
}

const blackrowDiv = document.getElementById("react-blackrow");
ReactDOM.render(<BlackRow />, blackrowDiv);
const whiterowDiv = document.getElementById("react-whiterow");
ReactDOM.render(<AnswersRow answers={[]} />, whiterowDiv);
const mycardsDiv = document.getElementById("react-mycards");
if (IS_PLAYER) {
    mycardsDiv.style.display = "flex";
}

const whiteRow = document.getElementById("whiterow");
// TODO: cambiare numero match
const socket = new WebSocket("ws://" + location.hostname + ":8080/ws?match=0");
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
    case "new_game":
        answers = [];
        ReactDOM.render(<BlackRow />, blackrowDiv);
        ReactDOM.render(<AnswersRow answers={[]}/>, whiterowDiv);
        if (IS_PLAYER) {
            const req = new XMLHttpRequest();
            req.addEventListener("load", () => {
                console.log(req.responseText);
                const resp = JSON.parse(req.responseText);
                const cards = resp.map(item => ({text: item.text, ID: item.Id}));
                console.log("My cards:", cards);
                ReactDOM.render(<MyCardsRow cards={cards} />, mycardsDiv);
            });
            req.open("GET", `/mycards?match_id=${0}`);
            req.send();
        }
        break;
    case "new_black":
        ReactDOM.render(<BlackRow card={<Card text={data.NewCard.text} black />}/>, blackrowDiv);
        break;
    case "new_white":
        cardText = getCardText(data)
        answers.push({ text: data.NewCard.text, total: 0, ID: data.NewCard.Id });
        ReactDOM.render(<AnswersRow answers={answers}/>, whiterowDiv);
        break;
    case "totals":
        totals = data.Totals
        for (const total of totals) {
            answers.find(a => a.ID == total.ID).total = total.Votes;
        }
        ReactDOM.render(<AnswersRow answers={answers} totals={totals}/>, whiterowDiv);
        break;
    default:
        alert("Unknown event " + eventName);
    }
}
socket.onclose = function () {
    console.log("Socket closed.");
}