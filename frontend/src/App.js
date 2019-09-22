import React from 'react'
import './App.css'
import BigText from './BigText'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBolt } from '@fortawesome/free-solid-svg-icons'
import axios from "axios"
import moment from "moment"

class App extends React.Component {
  state = {
    loading: true,
    location: "...",
    lastSeen: "...",
    lastSeenAbsolute: "...",
    battery: 0,
    charging: false
  }

  async componentDidMount() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/v1/location`)
    const at = moment(data.at)

    this.setState({
      loading: false,
      lastSeen: at.fromNow(),
      lastSeenAbsolute: at.format("dddd, MMMM Do YYYY, h:mm:ss a"),
      location: data.coarse_location,
      charging: false,
      battery: data.battery
    })
  }

  render() {
    if(this.state.loading) {
      return (
        <div className="App">
          Loading...
        </div>
      )
    }

    return (
      <div className="App">
        According to my phone, I was last seen
        <BigText title={this.state.lastSeenAbsolute}>{this.state.lastSeen}</BigText>
        somewhere in
        <BigText>
          <a href={`https://www.google.com/maps/search/${encodeURIComponent(this.state.location)}`} target="_blank">
            {this.state.location}.
          </a>
        </BigText>
        My battery was at
        <BigText>
          {this.state.battery}%
          {this.state.charging && <FontAwesomeIcon className="charging-icon" icon={faBolt} color="#0488f4" title="Charging"/>}
        </BigText>
        when this was reported.
      </div>
    )
  }
}

export default App
