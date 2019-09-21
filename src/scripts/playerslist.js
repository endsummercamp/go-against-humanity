import React, { Component } from "react";

export default class PlayersList extends Component {
	render() {
		const playerRows = this.props.players.map(player => <tr key={player.Id * 100 + player.Score}><td>{player.Username}</td><td>{player.Score}</td></tr>)
		const juryRows = this.props.jury.map(juror => <tr key={juror.Id}><td>{juror.Username}</td></tr>)
		return <div className="playersList">
			<h4 className="title is-4">Giocatori</h4>
			<table className="table">
				<tbody>{playerRows}</tbody>
			</table>
			<h4 className="title is-4">Pubblico</h4>
			<table className="table">
				<tbody>{juryRows}</tbody>
			</table>
		</div>
	}
}