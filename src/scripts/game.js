import React, { Component } from "react"
import ReactDOM from "react-dom"
import Timer from "./timer"
import Card from "./card"

if (!window.WebSocket) {
    alert("Your browser does not support WebSockets!")
}

const uiStateLabel = document.getElementById("ui-state-label");
function printUIState(text) {
	uiStateLabel.textContent = text;
}

class BlackRow extends Component {
    render() {
        return <div className="flex" id="blackrow">
            {("card" in this.props) ? this.props.card : <></>}
        </div>;
    }
}

let canVote = false;

class AnswersRow extends Component {
    tryVote(id) {
        if (IS_PLAYER) {
            alert("You're a player, you cannot vote!");
            return;
        }
        if (!canVote)
            return;
        const req = new XMLHttpRequest();
        req.open("PUT", `/matches/${MATCH_ID}/vote_card/${id}`);
        req.send();
        // canVote = false;
        return true;
    }
    render() {
        /* Expects:
           * a prop "answers", containing an array of {text, ID};
           * a prop "totals", containing an array of {Votes}.
         */
        let sum = 0;
        if (this.props.totals) {
            sum = this.props.totals.reduce((a, b) => a + b.Votes, 0);
		}

		const cards = this.props.answers.map((answer, i) => <Card text={answer.text} id={answer.ID} total={answer.total} sum={sum} onClick={(evt) => {
			console.log("Voted!");
			const success = this.tryVote(answer.ID);
			if (!success) return;
			evt.target.parentNode.classList.add("voted");
		}} key={i} />);

        return <div className="flex" id="blackrow">
            {cards}
        </div>;
    }
}

class MyCardsRow extends Component {
    submitCard(id) {
        if (!canPickCard) {
            alert("You cannot pick a card at this time!");
            return false;
        }
        const req = new XMLHttpRequest();
        req.open("PUT", `/matches/${MATCH_ID}/pick_card/${id}`);
        req.send();
        canPickCard = false;
        return true;
    }
    render() {
        /* Expects:
           * a prop "cards", containing an array of {text, ID};
         */
		const cards = this.props.cards.map((answer, i) => <Card text={answer.text} id={answer.ID} onClick={(evt) => {
			const success = this.submitCard(answer.ID);
			if (!success) return;
			evt.target.parentNode.classList.add("selected");
		}} key={i} />);
        return <>{cards}</>;
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

const socket = new WebSocket(`ws://${document.location.hostname}:8080/ws?match=${MATCH_ID}`);
printUIState("Connecting...");

socket.onopen = function() {
	printUIState("Waiting for a black card...");
    console.log("Opened socket.");
};

let answers = [];
let totals = [];
let canPickCard = false;

socket.onmessage = function (e) {
    const data = JSON.parse(e.data);
    console.log("Received", data);
	const { Name: eventName } = data;
    switch (eventName) {
    case "join_successful":
		// We joined successfully. Clear the UI.
        resetUI(data.SecondsUntilFinishPicking, data.InitialBlackCard.text);
        break;
	case "new_black":
		// A black card was chosen. Show it.
        mycardsDiv.style.display = "flex";
		ShowBlackCard(data.Duration, data.NewCard.text);
		if (IS_PLAYER)
			printUIState("Play your white card(s)!");
		else
			printUIState("Waiting for the players...");
        break;
	case "voting":
		// The voting phase has begun.
		if (IS_PLAYER)
			printUIState("The jurors are voting...");
		else
			printUIState("Vote for the best card!");
		canPickCard = false;
		// timerComponent.stop();
		mycardsDiv.style.display = "none";
		canVote = true;
		break;
	case "new_white":
		// A new white card (from the voting phase) was received.
        // let cardText = getCardText(data);
        answers.push({ text: data.NewCard.text, total: 0, ID: data.NewCard.Id });
        ReactDOM.render(<AnswersRow answers={answers}/>, whiterowDiv);
        break;
    case "vote_cast":
        totals = data.Totals;
        for (const total of totals) {
            answers.find(a => a.ID === total.ID).total = total.Votes;
        }
        ReactDOM.render(<AnswersRow answers={answers} totals={totals}/>, whiterowDiv);
        break;
    case "show_results":
		printUIState(""); // TODO
        resetUI();
        break;
    default:
        alert("Unknown event " + eventName);
    }
};

function resetUI(SecondsUntilFinishPicking, InitialBlackText) {
    for (const tag of document.getElementsByClassName("selected")) {
        tag.classList.remove("selected")
    }
    for (const tag of document.getElementsByClassName("voted")) {
        tag.classList.remove("voted")
    }
    answers = [];
    ReactDOM.render(<BlackRow />, blackrowDiv);
    ReactDOM.render(<AnswersRow answers={[]}/>, whiterowDiv);
    if (IS_PLAYER) {
        const req = new XMLHttpRequest();
        req.addEventListener("load", () => {
            const resp = JSON.parse(req.responseText);
            const cards = resp.map(item => ({text: item.text, ID: item.Id}));
            console.log("My cards:", cards);
            ReactDOM.render(<MyCardsRow cards={cards} />, mycardsDiv);
        });
        req.open("GET", `/mycards?match_id=${MATCH_ID}`);
        req.send();
    }
    if (SecondsUntilFinishPicking) {
        ShowBlackCard(SecondsUntilFinishPicking, InitialBlackText);
    }
}

socket.onclose = function () {
	printUIState("Lost connection!");
	// alert("Lost connection to the server.");
    console.log("Socket closed.");
};

const timer = document.getElementById("match-timer");

// Can be used to start a new game, or to "resume" an existing one
function ShowBlackCard(seconds_left, black_card_text) {
	canPickCard = true;
	ReactDOM.render(<Timer seconds={seconds_left} />, timer);
    ReactDOM.render(<BlackRow card={<Card text={black_card_text} black />}/>, blackrowDiv);
}

if (IS_ADMIN) {
	document.getElementById("admin-panel-new-blackcard")
		.addEventListener("click", () => {
			const req = new XMLHttpRequest();
			req.open("PUT", `/admin/matches/${MATCH_ID}/new_black_card`);
			req.send();
		});

	document.getElementById("admin-panel-end-voting")
		.addEventListener("click", () => {
			const req = new XMLHttpRequest();
			req.open("PUT", `/admin/matches/${MATCH_ID}/end_voting`);
			req.send();
		});
}