import React, { Fragment, Component } from 'react'
import $ from 'jquery'
import { trim, ltrim } from 'validator'
import autobind from 'autobind-decorator'
import axios from 'axios'

export class DbUserUpdate extends Component {

    constructor(props) {
        super(props)
        this.state = {
            prevUsername: this.props.user.username,
            id: this.props.user.id,
            username: this.props.user.username,
            password: this.props.user.password,
            superuser: this.props.user.superuser,
            usernameIsValid: false,
            passwordIsValid: false,
            formIsValid: false,
            usernameMessage: "",
            passwordMessage: "",
            alertMessage: ""
        }

        this.minUsernameLength = 4
        this.minPasswordLength = 4

    }

    @autobind
    showForm() {
        this.validateUsername()
        this.validatePassword()
        $("#update-dbuser-" + this.props.user.username).modal('show')
    }

    @autobind
    hideForm() {
        $("#update-dbuser-" + this.props.user.username).modal('hide')
    }

    @autobind
    updateUser() {
        axios.post('/api/v1/dbuser/update', {
                id: this.props.user.id,
                username: this.state.username,
                password: this.state.password,
                superuser: this.state.superuser
        }).then((res) => {
            if (res.data.error != null) {
                console.log("user update response: ", res.data)
                if (!res.data.error) {
                    this.hideForm()
                    this.props.hideCallback()
                    this.props.updateCallback()
                } else {
                    this.setState({
                        alertMessage: "Backend error"
                    })
                }
            }
        }).catch((err) => {
            this.setState({
                alertMessage: "Communication error"
            })
        })
    }

    @autobind
    onSubmit(event) {
        event.preventDefault()
        if (!this.state.formIsValid) {
            return
        }
        console.log("update user state:", this.state)
        this.updateUser()
    }

    validateUsername() {
        if (this.state.username.length > this.minUsernameLength) {
            this.setState({
                    usernameIsValid: true,
                    usernameMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    usernameIsValid: false,
                    usernameMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validatePassword() {
        if (this.state.password.length > this.minPasswordLength) {
            this.setState({
                    passwordIsValid: true,
                    passwordMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    passwordIsValid: false,
                    passwordMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validateForm() {
        if (this.state.usernameIsValid && this.state.passwordIsValid) {
            this.setState({ formIsValid: true })
        } else {
            this.setState({ formIsValid: false })
        }
    }

    @autobind
    onChangeUsername(event) {
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                username: newValue
            },
            () => { this.validateUsername() }
        )
    }

    @autobind
    onChangePassword(event) {
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                password: newValue
            },
            () => { this.validatePassword() }
        )
    }

    @autobind
    onIsAdmin(event) {
        this.setState({ superuser: event.target.checked })
    }

    @autobind
    componentDidMount() {
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
        return 'user-update-' + base + '-' + this.props.user.username
    }

    render() {
        return (
            <Fragment>

                <a onClick={this.showForm} className="dropdown-item">
                    <i className="far fa-edit"></i> Update
                </a>

                <div className="modal fade" id={"update-dbuser-" + this.props.user.username} tabIndex="-1" role="dialog"  ref="form">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">

                            <form acceptCharset="UTF-8" onSubmit={this.onSubmit} ref="form">

                                <div className="modal-header">
                                    <h5 className="modal-title">Update user [{this.props.user.username}]</h5>
                                    <button type="button" className="close" onClick={this.hideForm}>
                                        <span>&times;</span>
                                    </button>
                                </div>

                                <div className="modal-body">

                                    <div className="form-group">
                                        <label htmlFor={this.makeId("username")}>Username:</label>
                                        <input type="text" className="form-control" id={this.makeId("username")}
                                                    value={this.state.username}  onChange={this.onChangeUsername}/>
                                        <small className="form-text text-muted">{this.state.usernameMessage}</small>
                                    </div>

                                    <div className="form-group">
                                        <label htmlFor={this.makeId("password")}>Password:</label>
                                        <input id={this.makeId("password")} type="text" className="form-control"
                                                    value={this.state.password} onChange={this.onChangePassword}/>
                                        <small className="form-text text-muted">{this.state.passwordMessage}</small>
                                    </div>

                                    <div className="form-group form-check">
                                        <input id={this.makeId("superuser")} className="form-check-input"
                                                    type="checkbox" onChange={this.onIsAdmin} defaultChecked={this.props.user.superuser} />
                                        <label className="form-check-label" htmlFor={this.makeId("superuser")}> As admin</label>
                                    </div>

                                </div>

                                {this.showAlert()}

                                <div className="modal-footer">

                                    <button type="button" className="btn btn-sm btn-secondary" onClick={this.hideForm}>Close</button>
                                    <button type="submit" className={this.state.formIsValid ? "btn btn-sm btn-primary" : "btn btn-sm btn-primary disabled"}
                                                onClick={this.onSubmit}>Update</button>
                                </div>

                            </form>

                        </div>
                    </div>
                </div>
            </Fragment>
        )
    }
}

export default DbUserUpdate
