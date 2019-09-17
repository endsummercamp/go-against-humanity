import React, { PureComponent } from "react"

export default class Timer extends PureComponent {
	constructor(props) {
		super(props);
		this.state = {
			seconds_left: props.seconds
		};
		const interval = setInterval(() => {
			const seconds_left = this.state.seconds_left;
			if (seconds_left <= 0) {
				this.setState({
					seconds_left: 0
				});
				clearInterval(interval);
				return;
			}
			this.setState({
				seconds_left: seconds_left - 1
			});
		}, 1000);
	}

	render() {
		const seconds_left = this.state.seconds_left;
		const minutes = String(Math.floor(seconds_left / 60)).padStart(2, '0');
		const seconds = String(seconds_left % 60).padStart(2, '0');
		return <span>{minutes}:{seconds}</span>;
	}
}