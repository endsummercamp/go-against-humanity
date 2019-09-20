import React, { PureComponent, Component } from "react"
import {ProjectorViewBtn, NewBlackCardBtn, EndVotingBtn} from "./admin_buttons";

class UIStateLabel extends PureComponent {
	render() {
		return <div className="navbar-item">
			<strong id="ui-state-label">{this.props.text}</strong>
		</div>;
	}
}

export default class Navbar extends Component {
	render() {
		let adminButtons = <></>;
		if (IS_ADMIN)
			adminButtons = <div className="navbar-item">
				<div className="match-admin-panel">
					<ProjectorViewBtn />
					<NewBlackCardBtn />
					<EndVotingBtn />
				</div>
			</div>;
		
		return <nav className="navbar is-fixed-top">
			<a href="/" className="navbar-item">
				<img src="/public/img/ESC-logo-small.png" />
			</a>
			{window.isProjector && <a className="navbar-item">Projector view</a>}
			{adminButtons}
			<UIStateLabel text={this.props.uiStateText} />
		</nav>
	}
}