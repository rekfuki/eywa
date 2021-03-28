import React from "react";
import WebServiceBuilder from "./WebserviceBuilder";
import "./styles.css";

class SystemsBuilder extends React.Component {
  constructor(props) {
    super(props);
    let title = this.props.title != "" ? "Magento 2" : "";
    this.state = {
      title: title,
      code: "",
      description: "",
      protocol: "",
      isEdit: null,
      isDisplay: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.updateIsEdit = this.updateIsEdit.bind(this);
    this.openCloseGenerateToken = this.openCloseGenerateToken.bind(this);
  }

  updateIsEdit(e, value = null) {
    this.setState({
      isEdit: value
    });
  }

  handleChange(e) {
    this.setState({
      [e.target.name]: e.target.value
    });
  }

  openCloseGenerateToken(isDisplay) {
    this.setState({ isDisplay: isDisplay });
  }
  renderProtocol() {
    let protocols = [
      { key: "no_auth", label: "NO AUTH" },
      { key: "bearer_token", label: "Bearer Token" },
      { key: "basic_auth", label: "Basic Auth" }
    ];
    let protocolOptions = [];
    let loadedOption = this.state.protocol;

    protocols.forEach(function (m) {
      if (m.key == loadedOption) {
        protocolOptions.push(
          <option selected value={m.key}>
            {m.label}
          </option>
        );
      } else {
        protocolOptions.push(<option value={m.key}>{m.label}</option>);
      }
    });

    return (
      <div className="row">
        <div className="col-lg-4 p-r-0">
          <select
            value={this.state.protocol}
            name={"protocol"}
            className="form-control"
            onChange={(e) => this.handleChange(e, "protocol")}
          >
            {protocolOptions}
          </select>
        </div>
      </div>
    );
  }
  friendlyInput(name, placeholder = null, noEdit = false) {
    let showInput = false;
    if (this.state.isEdit == name) {
      showInput = true;
      if (noEdit && this.state[name]) {
        showInput = true;
      }
    }
    return (
      <div className="form-group header-text">
        {showInput ? (
          <input
            className="form-control"
            id={name}
            name={name}
            type="text"
            autoFocus={true}
            value={this.state[name]}
            onChange={this.handleChange}
            onBlur={this.updateIsEdit}
            placeholder={placeholder}
          />
        ) : (
          <div
            className="form-control"
            id={name}
            onClick={(e) => this.updateIsEdit(e, name)}
          >
            {this.state[name] ? (
              this.state[name]
            ) : (
              <div className="cl-placeholder">{placeholder}</div>
            )}
            <input type="hidden" name={name} value={this.state[name]} />
          </div>
        )}
      </div>
    );
  }

  render() {
    //console.log(this.state.title);
    return (
      <div className="container">
        <div className="col-lg-offset-2 col-lg-8">
          <div className="form-details">
            <form action="" method="POST">
              <br />
              <div className="testconnectionsection">
                <WebServiceBuilder
                  service_type="test"
                  data={this.state.isDisplay}
                  onChange={this.openCloseGenerateToken}
                />
              </div>
              {this.state.isDisplay ? (
                <div className="tokenkeysection">
                  <p className="tokenlabel">Dynamic Token</p>
                  <WebServiceBuilder service_type="token" />
                </div>
              ) : (
                ""
              )}
              <input type="submit" className="btn btn-default" value="Submit" />
            </form>
          </div>
        </div>
      </div>
    );
  }
}


export default SystemsBuilder;