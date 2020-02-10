import React, { Component, Fragment } from 'react'
import { Link } from 'react-router-dom'

import { store } from './main'


export class Menu extends Component {

    render() {
        return (
            <Fragment>
                <div className="mb-2">
                    <div className="dropdown-divider"></div>
                    <div className="dropdown-item active"><i className="fas fa-hammer"></i> Menu </div>

                    <Link className="dropdown-item" to="/users"><i className="fas fa-users"></i> App Users </Link>
                    <Link className="dropdown-item" to="/dbusers"><i className="fas fa-users-cog"></i> DB users </Link>
                    <Link className="dropdown-item" to="/dbs"><i className="fas fa-database"></i> Databases </Link>

                    <div className="dropdown-divider mb-3"></div>

                    <div className="dropdown-item" onClick={store.logout}><i className="fas fa-sign-out-alt"></i> Logout </div>
                    <div className="dropdown-divider mb-3"></div>
                </div>
            </Fragment>
        )
    }

}

export default Menu
