import React, { Fragment, Component } from 'react'
import $ from 'jquery'
import { trim, ltrim } from 'validator'
import autobind from 'autobind-decorator'
import axios from 'axios'

export class DbCreate extends Component {

    constructor(props) {
        super(props)

        this.initState = {
            dbname: "",
            owner: "",
            dbnameIsValid: false,
            ownerIsValid: false,
            formIsValid: false,
            dbnameMessage: "",
            ownerMessage: "",
            alertMessage: "",
            dbusers: []
        }
        this.state = this.initState

        this.minDbNameLength = 1
        this.minOwnerLength = 1
    }

    @autobind
    createDb() {
        axios.post('/api/v1/db/create', {
                dbname: this.state.dbname,
                owner: this.state.owner,
        }).then((res) => {
            if (res.data.error != null) {
                console.log("create: ", res.data)
                if (!res.data.error) {
                    this.hideForm()
                    this.clearForm()
                    this.props.updateCallback()
                } else {
                    this.setState({ alertMessage: "Backend error" + res.data.message  })
                }
            }
        }).catch((err) => {
                    this.setState({ alertMessage: "Communication error" })
        })
    }

    @autobind
    getDbUsers() {
        axios
            .get('/api/v1/dbuser/listall')
            .then((res) => {
                if (res.data.error != null) {
                    //console.log("get db users: ", res.data)
                    if (res.data.result == null) {
                        res.data.result = []
                    }
                    if (!res.data.error) {
                        this.setState({
                            dbusers: res.data.result,
                            alertMessage: ""
                        })
                    } else {
                        if (res.data.message == null) {
                            this.setState({
                                    alertMessage: "Backend error. "
                            })
                        } else {
                            this.setState({
                                    alertMessage: "Backend error." + res.data.message
                            })
                        }

                    }
                }
            }).catch((err) => {
                this.setState({
                    alertMessage: "Communication error"
                })
            })
    }

    @autobind
    showForm() {
        this.getDbUsers()
        $("#create-db").modal('show')
    }

    hideForm() {
        $("#create-db").modal('hide')
    }

    clearForm() {
        this.setState(this.initState)
    }

    @autobind
    onSubmit(event) {
        event.preventDefault()
        if (!this.state.formIsValid) {
            return
        }
        this.createDb()
        console.log("new db:", this.state)
    }

    validateDbName() {
        if (this.state.dbname.length > this.minDbNameLength) {
            this.setState({
                    dbnameIsValid: true,
                    dbnameMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    dbnameIsValid: false,
                    dbnameMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validateOwner() {
        if (this.state.owner.length > this.minOwnerLength) {
            this.setState({
                    ownerIsValid: true,
                    ownerMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    ownerIsValid: false,
                    ownerMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validateForm() {
        if (this.state.dbnameIsValid && this.state.ownerIsValid) {
            this.setState({ formIsValid: true })
        } else {
            this.setState({ formIsValid: false })
        }
    }


    @autobind
    onChangeDbName(event) {
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                dbname: newValue
            },
            () => { this.validateDbName() }
        )
    }

    @autobind
    onChangeOwner(event) {
        console.log(event.target.value)
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                owner: newValue
            },
            () => { this.validateOwner() }
        )
    }

    @autobind
    onIsAdmin(event) {
        this.setState({ superdb: event.target.checked })
    }

    showAlert() {
        if (this.state.alertMessage != "") {
            return (
                <div className="alert alert-warning border mx-4" role="alert">
                  <div className="text-center">{this.state.alertMessage}</div>
                </div>
            )
        }
    }

    @autobind
    makeId(base) {
        return 'db-create-' + base
    }

    @autobind
    renderDbUsers() {
        return this.state.dbusers.map((dbuser, index) => {
            return (
                <option key={index}>{dbuser.username}</option>
            )
        })
    }

    render() {
        return (
            <Fragment>

                <a onClick={this.showForm}>
                    <i className="fas fa-plus fa-lg"></i>
                </a>

                <div className="modal fade" id="create-db" tabIndex="-1" role="dialog"  ref="form">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">

                            <form acceptCharset="UTF-8" onSubmit={this.onSubmit} ref="form">

                                <div className="modal-header">
                                    <h5 className="modal-title">Create database</h5>
                                    <button type="button" className="close" onClick={this.hideForm}>
                                        <span>&times;</span>
                                    </button>
                                </div>

                                <div className="modal-body">

                                    <div className="form-group">
                                        <label htmlFor="dbname">Database name:</label>
                                        <input type="text" className="form-control" id="dbname" value={this.state.dbname}  onChange={this.onChangeDbName}/>
                                        <small className="form-text text-muted">{this.state.dbnameMessage}</small>
                                    </div>

                                    <div className="form-group">
                                        <label htmlFor="owner">Owner</label>
                                        <select className="form-control" id="owner" onChange={this.onChangeOwner}>
                                            <option></option>
                                            {this.renderDbUsers()}
                                        </select>
                                    </div>

                                </div>

                                {this.showAlert()}

                                <div className="modal-footer">
                                    <button type="button" className="btn btn-sm btn-secondary" onClick={this.hideForm}>Close</button>
                                    <button type="submit" className={this.state.formIsValid ? "btn btn-sm btn-primary" : "btn btn-sm btn-primary disabled"} onClick={this.onSubmit}>Create</button>
                                </div>

                            </form>

                        </div>
                    </div>
                </div>
            </Fragment>
        )
    }
}

export default DbCreate
