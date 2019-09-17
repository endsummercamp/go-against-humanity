import React, { PureComponent, Component } from "react"
import ReactDOM from "react-dom"
import Navbar from "./navbar"
import Card from "./card"

if (!window.WebSocket) {
    alert("Your browser does not support WebSockets!")
}

class BlackRow extends Component {
    render() {
        return <div className="flex" id="blackrow">
            {this.props.card || <></>}
        </div>;
    }
}

class WhiteRow extends Component {
	render() {
		return <div id="react-mycards">
			<div className="flex" id="blackrow"></div>
		</div>;
	}
}

let canPickCard = false,
	canVote = false;

class MyCardsRow extends Component {
	constructor(props) {
		super(props);
		this.state = {
			selectedCard: null
		};
	}

	submitCard(id) {
        if (!canPickCard) {
            // alert("You cannot pick a card at this time!");
            return false;
        }
        const req = new XMLHttpRequest();
        req.open("PUT", `/matches/${MATCH_ID}/pick_card/${id}`);
        req.send();
        canPickCard = false;
        return true;
	}

	render() {
		const cards = this.props.cards.map(answer => <Card
			text={answer.text}
			id={answer.ID}
			selected={answer.ID == this.state.selectedCard}
			onClick={() => {
				const success = this.submitCard(answer.ID);
				if (!success)
					return;
				this.setState({selectedCard: answer.ID});
			}}
			key={answer.ID} />
		);
        return <div id="react-mycards">
			{cards}
		</div>
	}
}

class AnswersRow extends Component {
	constructor(props) {
		super(props);
		this.state = {
			votedCard: null
		};
	}

	tryVote(id) {
        if (IS_PLAYER) {
            // alert("You're a player, you cannot vote!");
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
        // Expects:
        // * a prop "answers", containing an array of {text, ID};
        // * a prop "totals", containing an array of {Votes}.
        let sum = 0;
        if (this.props.totals) {
            sum = this.props.totals.reduce((a, b) => a + b.Votes, 0);
		}

		const cards = this.props.answers.map(answer => <Card
			voted={answer.ID == this.state.votedCard}
			text={answer.text}
			id={answer.ID}
			total={answer.total}
			sum={sum}
			onClick={(evt) => {
				console.log("Voted!");
				const success = this.tryVote(answer.ID);
				if (!success) return;
				this.setState({votedCard: answer.ID});
			}}
			key={answer.ID} />);

        return <div className="flex" id="blackrow">
            {cards}
        </div>;
    }
}

class Game extends Component {
	constructor(props) {
		super(props);
		this.socket = new WebSocket(`ws://${document.location.hostname}:8080/ws?match=${MATCH_ID}`);
		this.state = {
			// Navbar state
			timerState: {
				enabled: false,
				seconds: 0
			},
			uiStateText: "Connecting...",
			// Game UI state
			blackCard: null,
			myCards: [],
			answers: [],
		};
		this.socket.onopen = () => {
			this.setState(Object.assign(this.state, {uiStateText: "Waiting for a black card..."}));
			// console.log("Opened socket.");
		};
		this.socket.onmessage = e => {
			const data = JSON.parse(e.data);
			console.log("Received", data);
			const eventName = data.Name;
			switch (eventName) {
			case "join_successful":
				// We joined successfully. Clear the UI.
				this.resetUI();
				if (data.SecondsUntilFinishPicking)
					this.showBlackCard(data.SecondsUntilFinishPicking, data.InitialBlackText.text);	
				break;
			case "new_black":
				// A black card was chosen. Show it.
				// mycardsDiv.style.display = "flex";
				this.showBlackCard(data.Duration, data.NewCard.text);
				break;
			case "voting":
				// The voting phase has begun.
				canPickCard = false;
				canVote = true;
				this.setState(Object.assign(this.state, {
					timerState: {
						enabled: false,
						seconds: 0
					},
					uiStateText: IS_PLAYER
						? "The jurors are voting..."
						: "Vote for the best card!",
					myCards: []
				}))
				break;
			case "new_white":
				// A new white card (from the voting phase) was received.
				this.setState(Object.assign(this.state, {
					answers: this.state.answers.concat({ text: data.NewCard.text, total: 0, ID: data.NewCard.Id })
				}));
				break;
			case "vote_cast":
				let totals = data.Totals;
				let answers = this.state.answers;
				for (const total of totals) {
					answers.find(a => a.ID === total.ID).total = total.Votes;
				}
				this.setState(Object.assign(this.state, {
					answers,
					totals,
				}));
				break;
			case "show_results":
				this.setState(Object.assign(this.state, {
					uiStateText: "", // TODO
				}));
				resetUI();
				break;
			default:
				alert("Unknown event " + eventName);
			}
		};
	}

	showBlackCard(duration, text) {
		canPickCard = true;
		console.log("showBlackCard:", duration);
		this.setState(Object.assign(this.state, {
			timerState: {
				enabled: true,
				seconds: duration
			},
			uiStateText: IS_PLAYER
				? "Play your white card!"
				: "Waiting for the players...",
			blackCard: <Card text={text} black />
		}));
	}

	resetUI() {
		/* todo
		for (const tag of document.getElementsByClassName("selected")) {
			tag.classList.remove("selected")
		}
		for (const tag of document.getElementsByClassName("voted")) {
			tag.classList.remove("voted")
		}
		*/
		this.setState(Object.assign(this.state, {
			blackCard: null,
			answers: []
		}));
		if (IS_PLAYER)
			this.fetchMyCards();
	}

	fetchMyCards() {
		if (!IS_PLAYER)
			return;
		const req = new XMLHttpRequest();
		req.addEventListener("load", () => {
			const resp = JSON.parse(req.responseText);
			const cards = resp.map(item => ({text: item.text, ID: item.Id}));
			console.log("My cards:", cards);
			this.setState(Object.assign(this.state, {myCards: cards}));
		});
		req.open("GET", `/mycards?match_id=${MATCH_ID}`);
		req.send();
	}

	render() {
		return <>
			<Navbar timerState={this.state.timerState} uiStateText={this.state.uiStateText} />
			<BlackRow card={this.state.blackCard} />
			<MyCardsRow cards={this.state.myCards} />
			<AnswersRow answers={this.state.answers} />
		</>;
	}
}

ReactDOM.render(<Game />, document.getElementById("react-game"));