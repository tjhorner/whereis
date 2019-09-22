import React from 'react'
import moment from "moment"

class FromNow extends React.Component {
  state = {
    interval: null
  }

  componentDidMount() {
    let interval = setInterval(() => this.forceUpdate(), 1000)
    this.setState({ interval })
  }

  componentWillUnmount() {
    clearInterval(this.state.interval)
  }

  render() {
    return <span>{moment(this.props.date).fromNow()}</span>
  }
}

export default FromNow
