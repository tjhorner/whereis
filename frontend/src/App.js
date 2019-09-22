import React from 'react'
import { BrowserRouter as Router, Route, Switch } from "react-router-dom"
import AnonPage from './AnonPage'
import KeyPage from './KeyPage'
import Error from './Error';
import NotFoundPage from './NotFoundPage';

class App extends React.Component {
  render() {
    return (
      <Router>
        <Switch>
          <Route path="/" exact component={AnonPage}/>
          <Route path="/:id" exact component={KeyPage}/>
          <Route component={NotFoundPage}/>
        </Switch>
      </Router>
    )
  }
}

export default App
