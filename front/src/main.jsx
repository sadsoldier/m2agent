
import './app.scss'
import 'bootstrap'

import React from 'react'
import ReactDOM from 'react-dom'

import { Router, Route, Switch } from 'react-router-dom'
import { createBrowserHistory } from 'history'

import { Provider } from 'mobx-react';
import * as Cookies from 'js-cookie'

import { decorate, observable, action } from "mobx"

import NotFound from './NotFound'
import NotEnough from './NotEnough'

import Login from './Login'
import Users from './Users'


import DbUsers from './DbUsers'
import Dbs from './Dbs'


export const history = createBrowserHistory()

const cookieName = "session"

class Store {
    @observable username
    @observable limit
    @observable isadmin

    constructor() {
        this.username = ""
        this.isadmin = false
        this.fileLimit = 10
        this.filePattern = "*"
        this.bucketLimit = 10
        this.bucketPattern = "*"
    }

    @action login = (user) => {
        this.username = user.username
        this.isadmin = user.isadmin
    }
    @action logout = () => {
        this.username = ""
        this.isadmin = false
        Cookies.remove(cookieName)
        history.push("/login")

    }
}

export const store = new Store()

export function checkLogin(level) {
    let cookie = Cookies.get(cookieName)
    if (store.username == "" || cookie == null) {
        history.push("/login")
    }
    if (level == "admin" && store.isadmin == false) {
        //history.push("/login")
        history.goBack()
    }
}


ReactDOM.render(
    <Router history={history}>
        <Switch>
            <Route exact path="/login" component={Login} />
            <Route exact path="/users" component={Users} />
            <Route exact path="/dbusers" component={DbUsers} />
            <Route exact path="/dbs" component={Dbs} />
            <Route exact path="/" component={Dbs} />

            <Route exact path="/notenough" component={NotEnough} />

            <Route path="*" component={NotFound} />
        </Switch>
    </Router>,
    document.getElementById('root')
)
