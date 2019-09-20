import React, { PureComponent, Component } from "react"

class NewBlackCardBtn extends PureComponent {
	render() {
		return <button className="button" id="admin-panel-new-blackcard" onClick={() => {
			const req = new XMLHttpRequest();
			req.open("PUT", `/admin/matches/${MATCH_ID}/new_black_card`);
			req.send();
		}}>New black card</button>;
	}
}

class ProjectorViewBtn extends PureComponent {
	render() {
		return <button className="button" id="admin-panel-projector-view" onClick={() => {
			window.isProjector ^= 1;
		}}>Toggle projector view</button>;
	}
}

class EndVotingBtn extends PureComponent {
	render() {
		return <button className="button" id="admin-panel-end-voting" onClick={() => {
			const req = new XMLHttpRequest();
			req.open("PUT", `/admin/matches/${MATCH_ID}/end_voting`);
			req.send();
		}}>End voting</button>;
	}
}

export {NewBlackCardBtn, ProjectorViewBtn, EndVotingBtn};