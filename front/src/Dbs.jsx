import React, { Component, Fragment } from 'react'
import autobind from 'autobind-decorator'
import axios from 'axios'
import dayjs from 'dayjs'

import { clone } from 'lodash'
import { formatNumber, fileSize } from 'humanize-plus'

import { Layout } from './Layout'
import { checkLogin } from './main'

import { DbCreate } from './DbCreate'
import { DbOption } from './DbOption'

import { Pager } from './Pager'

export class Dbs extends Component {

    constructor(props) {
        super(props)

        this.state = {
            dbs: [],
            offset: 0,
            limit: 5,
            total: 0,
            alertMessage: "",
            dbusers: []
        }
    }

    @autobind
    listDbs() {
        axios.post('/api/v1/db/list', {
                limit: this.state.limit,
                offset: this.state.offset
        }).then((res) => {
            if (res.data.error != null) {
                console.log("list dbs response: ", res.data)
                if (res.data.result.dbs == null) {
                    res.data.result.dbs = []
                }
                if (!res.data.error) {
                    this.setState({
                        dbs: res.data.result.dbs,
                        total: res.data.result.total,
                        offset: res.data.result.offset,
                        limit: res.data.result.limit,
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
    showAlert() {
        if (this.state.alertMessage.length > 0) {
            return (
                <div className="alert alert-warning border mx-4" role="alert">
                  <div className="text-center">{this.state.alertMessage}</div>
                </div>
            )
        }
    }

    @autobind
    changeOffset(newOffset) {
        this.setState({ offset: newOffset }, () => { this.listDbs() })
    }

    @autobind
    onChangeLimit(event) {
        const newLimit = Number(event.target.value)
        var newOffset = Math.floor(this.state.offset / newLimit) * newLimit
        this.setState({ limit: newLimit, offset: newOffset }, () => { this.listDbs() })
    }

    @autobind
    renderTable() {
        return this.state.dbs.map((db, index) => {
            const { dbname, owner, size, numbackends } = db
            const theDb = db
            return (
                <tr key={index}>
                    <td>{index + this.state.offset + 1}</td>
                    <td><DbOption db={theDb} updateCallback={this.listDbs}>{dbname}</DbOption></td>
                    <td>{owner}</td>
                    <td>{fileSize(size, 0)}</td>
                    <td>{numbackends}</td>
                </tr>
            )
        })
    }

    render() {
        return (
            <Fragment>
                <Layout>
                    <div className="container-fluid">
                        <div className="row justify-content-between mb-3">

                            <h5>
                                <i className="fas fa-database fa-sm"></i>
                                <span> Databases </span>
                                <a onClick={this.listDbs}><i className="fas fa-sync fa-sm"></i></a>
                                <small> {dayjs().format('hh:mm:ss')}</small>
                            </h5>
                            <div className="ml-auto">
                                <DbCreate updateCallback={this.listDbs}/>
                            </div>

                        </div>
                    </div>

                    {this.showAlert()}

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
                                <th>name</th>
                                <th>owner</th>
                                <th>size</th>
                                <th>#c</th>
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
        this.listDbs()
    }
}

export default Dbs
