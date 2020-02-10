import React, { Component, Fragment } from 'react'
import autobind from 'autobind-decorator'
import axios from 'axios'
import dayjs from 'dayjs'

import { clone } from 'lodash'

import { Layout } from './Layout'
import { checkLogin } from './main'

import { DbUserCreate } from './DbUserCreate'
import { DbUserOption } from './DbUserOption'


import { Pager } from './Pager'

export class DbUsers extends Component {

    constructor(props) {
        super(props)

        this.state = {
            users: [],
            offset: 0,
            limit: 5,
            total: 0,
            alertMessage: ""
        }
    }

    @autobind
    listUsers() {
        axios.post('/api/v1/dbuser/list', {
                limit: this.state.limit,
                offset: this.state.offset
        }).then((res) => {
            if (res.data.error != null) {
                console.log("list users response: ", res.data)
                if (res.data.result.users == null) {
                    res.data.result.users = []
                }
                if (!res.data.error) {
                    this.setState({
                        users: res.data.result.users,
                        total: res.data.result.total,
                        offset: res.data.result.offset,
                        limit: res.data.result.limit
                    })
                } else {
                    if (res.data.message == null) {
                        this.setState({
                                alertMessage: "Backend error"
                        })
                    } else {
                        this.setState({
                                alertMessage: res.data.message
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
    changeOffset(newOffset) {
        this.setState({ offset: newOffset }, () => { this.listUsers() })
    }

    @autobind
    onChangeLimit(event) {
        const newLimit = Number(event.target.value)
        var newOffset = Math.floor(this.state.offset / newLimit) * newLimit
        this.setState({ limit: newLimit, offset: newOffset }, () => { this.listUsers() })
    }

    @autobind
    renderTable() {
        return this.state.users.map((user, index) => {
            const { username, superuser } = user
            const theUser = user
            return (
                <tr key={index}>
                    <td>{index + this.state.offset + 1}</td>
                    <td><DbUserOption user={theUser} updateCallback={this.listUsers}>{username}</DbUserOption></td>
                    <td>{this.renderIsAdmin(superuser)}</td>
                </tr>
            )
        })
    }

    @autobind
    renderIsAdmin(superuser) {
        if (superuser) {
            return <i className="fas fa-user-plus"></i>;
        } else {
            return <i className="fas fa-microchip"></i>;
        }
    }

    render() {
        return (
            <Fragment>
                <Layout>
                    <div className="container-fluid">
                        <div className="row justify-content-between mb-3">
                            <h5>
                                <i className="fas fa-users-cog"></i>
                                <span> DB users </span>
                                <a onClick={this.listUsers}>
                                    <i className="fas fa-sync fa-sm"></i>
                                </a>
                                <small> {dayjs().format('hh:mm:ss')}</small>
                            </h5>
                            <div className="ml-auto">
                                <DbUserCreate updateCallback={this.listUsers}/>
                            </div>
                        </div>
                    </div>

                    <div className="d-inline-flex mb-1">
                        <select className="custom-select" id="page-limit" value={this.state.limit} onChange={this.onChangeLimit}>
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="25">25</option>
                            <option value="50">50</option>
                        </select>
                    </div>

                    <table className="table table-striped table-hover table-sm">

                        <thead className="thead-light">
                            <tr>
                                <th>#</th>
                                <th>login name</th>
                                <th>mode</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.renderTable()}
                        </tbody>
                    </table>

                    <Pager total={this.state.total} limit={this.state.limit} offset={this.state.offset} callback={this.changeOffset} />

                </Layout>
            </Fragment>
        )
    }

    componentDidMount() {
        checkLogin()
        this.listUsers()
    }
}

export default DbUsers
