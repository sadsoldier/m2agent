import React, { Fragment, Component } from 'react'
import $ from 'jquery'
import autobind from 'autobind-decorator'
import axios from 'axios'

import { DbUpdate } from './DbUpdate'
import { DbDelete } from './DbDelete'

export class DbOption extends Component {

    //constructor(props) {
        //super(props)
    //}

    @autobind
    showForm() {
        $("#option-db-" + this.props.db.dbname).modal('show')
    }

    @autobind
    hideForm() {
        $("#option-db-" + this.props.db.dbname).modal('hide')
    }

    render() {
        return (
            <Fragment>

                <a onClick={this.showForm}>{this.props.children}</a>

                <div className="modal fade bd-example-modal-sm" id={"option-db-" + this.props.db.dbname}
                                                                    tabIndex="-1" role="dialog"  ref="form">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">

                                <div className="modal-header">
                                    <h5 className="modal-title">Database [{this.props.db.dbname}] option</h5>
                                    <button type="button" className="close"
                                                                onClick={this.hideForm}>&times;</button>
                                </div>

                                <div className="modal-body">
                                            <DbUpdate db={this.props.db}
                                                            updateCallback={this.props.updateCallback}
                                                            hideCallback={this.hideForm} />
                                            <div className="dropdown-divider"></div>
                                            <DbDelete db={this.props.db}
                                                            updateCallback={this.props.updateCallback}
                                                            hideCallback={this.hideForm} />

                                </div>

                                <div className="modal-footer">
                                    <button type="button" className="btn btn-sm btn-secondary"
                                                            onClick={this.hideForm}>Close</button>
                                </div>

                        </div>
                    </div>
                </div>
            </Fragment>
        )
    }
}

export default DbOption
