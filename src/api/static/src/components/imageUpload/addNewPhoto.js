import {Component} from "react";
import React from "react";
import PropTypes from "prop-types";
import ImageUpload from "../imageUpload/index";

export default class AddNewPhoto extends Component {

    render() {
        return (
            <ImageUpload
                onError={this.onError}
                onLoad={this.onLoad}
                onAdd={this.props.onAdd}
                height={this.props.height}
                disabled={this.props.disabled}
                header={this.props.header}
                folder={this.props.folder}
            />
        )
    }
}

AddNewPhoto.propTypes = {
    onAdd: PropTypes.func.isRequired
};
