import React, { PureComponent, Component } from "react"
import {NewBlackCardBtn, EndVotingBtn} from "./admin_buttons";
import Timer from "./timer"

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
					<NewBlackCardBtn />
					<EndVotingBtn />
				</div>
			</div>;
		
		return <nav className="navbar is-fixed-top">
			<a className="navbar-item">
				<img src="/public/img/ESC-logo-small.png" />
			</a>
			<Timer enabled={this.props.timerState.enabled} seconds={this.props.timerState.seconds} />
			{adminButtons}
			<UIStateLabel text={this.props.uiStateText} />
		</nav>
	}
}