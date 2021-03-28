import React from "react";
import { Tabs, Tab } from "react-bootstrap";

class WebServiceBuilder extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      method: "get",
      url: "",
      params: [{ key: "", value: "", label: "" }],
      header: [{ key: "", value: "", label: "" }],
      body: [{ key: "", value: "", label: "" }],
      body_raw: "",
      suggestionOptions: [],
      isVisible: false,
      currentElementText: []
    };
    this.handleRowKeyUp = this.handleRowKeyUp.bind(this);
    this.handleRowChange = this.handleRowChange.bind(this);
    this.handleBulkChange = this.handleBulkChange.bind(this);
    this.removeClick = this.removeClick.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.changeEditView = this.changeEditView.bind(this);
    this.resetEditView = this.resetEditView.bind(this);
    this.generateBulkEditValue = this.generateBulkEditValue.bind(this);
    this.handleClick = this.handleClick.bind(this);
    this.removeSuggesstions = this.removeSuggesstions.bind(this);
  }

  handleChange(e, name = null) {
    if (!name) {
      name = e.target.name;
    }
    this.setState({
      [name]: e.target.value
    });
  }

  changeEditView() {
    this.setState({
      isVisible: !this.state.isVisible
    });
  }

  resetEditView() {
    this.setState({
      isVisible: false
    });
  }

  handleRowKeyUp(e, value) {
    let rowValue = this.state[value];
    let len = rowValue.length;
    if (
      (len > 0 &&
        (rowValue[len - 1]["key"].length > 0 ||
          rowValue[len - 1]["value"].length > 0)) ||
      rowValue[len - 1]["label"].length > 0
    ) {
      this.setState((prevState) => ({
        [value]: [...prevState[value], { key: "", value: "", label: "" }]
      }));
    }
  }
  getSuggesstionOptions(value) {
    return ["DYNAMIC_TOKEN"].filter(
      (item) => item.toLowerCase().indexOf(value.toLowerCase()) !== -1
    );
  }
  handleClick(e, i, value) {
    let rowValues = this.state[value];
    if (this.state.currentElementText["name"] == "value") {
      rowValues[this.state.currentElementText["i"]][
        this.state.currentElementText["name"]
      ] = this.subStrValue + e.target.innerText;
      this.setState({
        [value]: rowValues,
        value,
        suggestionOptions: []
      });
    }
    this.isVisibleGenerateToken();
  }
  handleRowChange(e, i, value) {
    clearTimeout(this.timeout);
    this.subStrValue = "";
    let rowValues = this.state[value];
    rowValues[i][e.target.name] = e.target.value;
    let currentElementText = [];
    currentElementText["i"] = i;
    currentElementText["name"] = e.target.name;
    this.setState({
      [value]: rowValues,
      value,
      currentElementText: currentElementText
    });
    if (e.target.name == "value") {
      let value = e.target.value;
      let substr = value.split(" ");
      if (substr.length > 1) {
        value = substr[substr.length - 1];
      } else {
        value = substr[0];
      }
      for (var i = 0; i <= substr.length - 2; i++) {
        this.subStrValue += substr[i] + " ";
      }
      if (value.length > 0) {
        this.timeout = setTimeout(() => {
          const suggestionOptions = this.getSuggesstionOptions(value);
          this.setState({
            suggestionOptions
          });
        }, 250);
      } else {
        this.setState({
          suggestionOptions: []
        });
      }
      this.isVisibleGenerateToken();
    }
  }

  handleBulkChange(e, name) {
    let rowValues = [];
    let rowValue = this.state[name];
    let rows = e.target.value.split(/\r?\n/);
    if (rows.length) {
      rows.forEach(function (rowItem) {
        let row = rowItem.split("::");
        if (row.length == 3) {
          if (row[0] && row[1] && row[2]) {
            rowValues.push({ label: row[0], key: row[1], value: row[2] });
          }
        }
      });
      rowValues.push({ key: "", value: "", label: "" });
      this.isVisibleGenerateToken();
      this.setState((prevState) => ({
        [name]: rowValues
      }));
    }
  }

  isVisibleGenerateToken() {
    let flag = 0;
    let names = ["params", "header", "body"];
    let self = this;

    if (names.length > 0) {
      names.forEach(function (name) {
        self.state[name].forEach(function (rowItem, i) {
          let rowItemValue = self.state[name][i]["value"];
          let n = rowItemValue.indexOf("DYNAMIC_TOKEN");
          if (n >= 0) {
            flag = 1;
          }
        });
      });
    }
    this.props.onChange(flag);
  }

  removeClick(e, i, name) {
    let rows = this.state[name];
    rows.splice(i, 1);
    this.setState({ [name]: rows });
    this.isVisibleGenerateToken(rows);
  }
  removeSuggesstions() {
    this.setState({ open: !this.state.open });
  }
  renderKeyValuePair(name) {
    let suggestion = this.state.suggestionOptions.map(
      function (select, i) {
        return (
          <li
            key={i}
            className="list-group-item"
            name={"suggestions"}
            value={select}
            onClick={(e) => this.handleClick(e, i, name)}
          >
            {select}
          </li>
        );
      }.bind(this)
    );

    let keyValuePair = this.state[name].map((el, i) => (
      <tr key={i}>
        <td>
          <div className="form-group">
            <input
              placeholder={"Label"}
              name={"label"}
              className="form-control"
              value={el.label || ""}
              onChange={(e) => this.handleRowChange(e, i, name)}
              onKeyUp={(e) => this.handleRowKeyUp(e, name)}
              id={"label"}
              type="text"
            />
          </div>
        </td>
        <td>
          <div className="form-group">
            <input
              placeholder={"Key"}
              name={"key"}
              className="form-control"
              value={el.key || ""}
              onChange={(e) => this.handleRowChange(e, i, name)}
              onKeyUp={(e) => this.handleRowKeyUp(e, name)}
              id={"key"}
              type="text"
            />
          </div>
        </td>
        <td>
          <div className="form-group">
            <input
              placeholder={"Value"}
              name={"value"}
              className="form-control"
              value={el.value || ""}
              onChange={(e) => this.handleRowChange(e, i, name)}
              onKeyUp={(e) => this.handleRowKeyUp(e, name)}
              onBlur={this.removeSuggesstions}
              id={name + "[" + i + "]"}
              type="text"
            />
            {i === this.state.currentElementText["i"] ? (
              <div className="suggestionValues">{suggestion}</div>
            ) : (
              ""
            )}
          </div>
        </td>
        <td>
          <button
            type="button"
            className="btn btn-sm"
            onClick={(e) => this.removeClick(e, i, name)}
          >
            <i className="fa fa-times" aria-hidden="true" />
          </button>
        </td>
      </tr>
    ));

    return (
      <table className="table table-striped" id="dataTable">
        <thead>
          <tr>
            <th className="text-center">Label</th>
            <th className="text-center">Key</th>
            <th className="text-center">Value</th>
            <th className="text-center" />
          </tr>
        </thead>
        <tbody id={name}>{keyValuePair}</tbody>
      </table>
    );
  }

  renderUrlBuilder() {
    let methods = [
      { key: "get", label: "GET" },
      { key: "post", label: "POST" },
      { key: "delete", label: "DELETE" },
      { key: "patch", label: "PATCH" },
      { key: "put", label: "PUT" }
    ];
    let options = [];
    let loadedOption = this.state.method;

    methods.forEach(function (m) {
      if (m.key == loadedOption) {
        options.push(
          <option selected value={m.key}>
            {m.label}
          </option>
        );
      } else {
        options.push(<option value={m.key}>{m.label}</option>);
      }
    });

    return (
      <div className="row">
        <div className="col-lg-2 p-r-0">
          <select
            value={this.state.method}
            name={this.props.service_type + "_method"}
            className="form-control"
            onChange={(e) => this.handleChange(e, "method")}
          >
            {options}
          </select>
        </div>
        <div className="col-lg-10 p-l-0">
          <div className="form-group">
            <input
              placeholder={"Url"}
              name={this.props.service_type + "_url"}
              className="form-control"
              value={this.state.url}
              onChange={(e) => this.handleChange(e, "url")}
              id={this.props.service_type + "_url"}
              type="text"
            />
          </div>
        </div>
      </div>
    );
  }

  generateBulkEditValue(name) {
    let bulkValue = "";
    this.state[name].map(function (row) {
      if (row["label"] && row["key"] && row["value"]) {
        bulkValue +=
          row["label"] + "::" + row["key"] + "::" + row["value"] + "\n";
      }
    });

    return (
      <div className="form-group ">
        <textarea
          className="form-control"
          name={this.props.service_type + "_" + name}
          onChange={(e) => this.handleBulkChange(e, name)}
          rows="10"
        >
          {bulkValue}
        </textarea>
        <input
          type="hidden"
          name={this.props.service_type + "_" + name + "_json"}
          value={JSON.stringify(this.state[name])}
        />
      </div>
    );
  }

  renderTab(name) {
    return (
      <div className="tebselection">
        {this.state.isVisible ? (
          <p className="title" onClick={this.changeEditView}>
            Key Value Edit
          </p>
        ) : (
          <p className="title" onClick={this.changeEditView}>
            Bulk Edit
          </p>
        )}
        {this.state.isVisible ? (
          <span>{this.generateBulkEditValue(name)}</span>
        ) : (
          <span>{this.renderKeyValuePair(name)}</span>
        )}
      </div>
    );
  }
  render() {
    return (
      <div className="testonnectionsection">
        {this.renderUrlBuilder()}
        <Tabs defaultActiveKey={1} onSelect={this.resetEditView}>
          <Tab eventKey={1} title="Params">
            {this.renderTab("params")}
          </Tab>
          <Tab eventKey={2} title="Header">
            {this.renderTab("header")}
          </Tab>
          <Tab eventKey={3} title="Body">
            <Tabs>
              <Tab.Pane eventKey={3.1} title="Form Data">
                {this.renderTab("body")}
              </Tab.Pane>
              <Tab.Pane eventKey={3.2} title="Raw">
                <div className="form-group">
                  <textarea
                    className="form-control"
                    id={this.props.service_type + "_body_raw"}
                    rows="10"
                    placeholder="Provide raw information"
                    value={this.state.body_raw}
                    onChange={(e) => this.handleChange(e, "body_raw")}
                  />
                </div>
              </Tab.Pane>
            </Tabs>
          </Tab>
        </Tabs>
      </div>
    );
  }
}
export default WebServiceBuilder;
