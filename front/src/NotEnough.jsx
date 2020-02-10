import React, { Fragment, Component } from 'react'

import { Layout } from './Layout'
//import { checkLogin } from './main'

export class NotEnough extends Component {

    render() {
        return (
            <Fragment>
                <Layout>
                    <div className="row">
                        <h4>401 Not Enough Rights</h4>
                    </div>
                </Layout>
            </Fragment>
        )
    }

    componentDidMount() {
        //checkLogin("")
    }

}

export default NotEnough
