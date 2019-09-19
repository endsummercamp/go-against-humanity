import React, { Component } from "react"

function dashFix(content){
    return content.replace(/_/g, '<div class="long-dash"></div>')
}

export default class Card extends Component {
    render() {
        // https://stackoverflow.com/a/6040258
        // https://stackoverflow.com/a/8541575
        let style = {};
        if (this.props.total) {
            const percentage = this.props.total / this.props.sum;
            style.height = percentage * 100 + "%";
		}
		const classes = [
			"card",
			this.props.black ? "card-black" : "card-white"
		];
		if (this.props.selected)
			classes.push("selected");
        if (this.props.text.length > 40)
            classes.push("small-text");

        let classes_top = "card-top card-content";
        if (this.props.total)
            classes_top += " with-z-index-fix";

        let classes_middle = "card-middle";
        if (this.props.total)
            classes_middle += " active";

		return <div className={classes.join(" ")} onClick={this.props.onClick}>
            <div className={classes_top} dangerouslySetInnerHTML={{__html:dashFix(this.props.text)}}></div>
            <div className={classes_middle}>
                <div className="card-votes">{this.props.total || ""}</div>
            </div>
            <div className="card-bottom">
                Cards Against Humanity
            </div>
            <div style={{position: "absolute", top: 0, left: 0, bottom: 0, right: 0, zIndex: 0}}>
                <div style={style} className="vote-bg"></div>
            </div>
        </div>
    }
}