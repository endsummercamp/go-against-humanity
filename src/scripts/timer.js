import React, { PureComponent } from "react"

export default class Timer extends PureComponent {
	constructor(props) {
		super(props);
		this.state = {
			now: Math.round((new Date()).getTime()/1000)
		};

		setInterval(() => {
			this.setState({
				now: Math.round((new Date()).getTime()/1000)
			});
		}, 1000);
	}

	render() {
		const now = this.state.now;
		let minutes = 0;
		let seconds = 0;

		if(this.props.expires > now) {
			const diff = this.props.expires - now;
			minutes = Math.round(diff/60);
			seconds = diff - minutes * 60;
		}
		minutes = ('0' + minutes).slice(-2);
		seconds = ('0' + seconds).slice(-2);
		return <a className="navbar-item">{minutes}:{seconds}</a>;
	}
}