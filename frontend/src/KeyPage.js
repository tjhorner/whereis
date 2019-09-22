import React from 'react'
import './App.css'
import BigText from './BigText'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBolt } from '@fortawesome/free-solid-svg-icons'
import axios from "axios"
import moment from "moment"
import Error from './Error'
import { Map, TileLayer, Marker, Circle } from "react-leaflet"

class KeyPage extends React.Component {
  state = {
    loading: true,
    location: "...",
    lat: 0,
    lng: 0,
    accuracy: 0,
    lastSeen: "...",
    lastSeenAbsolute: "...",
    battery: 0,
    charging: false,
    interval: null,
    error: null
  }

  async updateData() {
    try {
      const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/v1/location`, {
        params: { key: this.props.match.params.id }
      })

      if(data.error) {
        this.setState({ error: data.error })
        return
      }

      const at = moment(data.at)

      this.setState({
        loading: false,
        lastSeen: at.fromNow(),
        lastSeenAbsolute: at.format("dddd, MMMM Do YYYY, h:mm:ss a"),
        location: data.coarse_location,
        lat: data.latitude,
        lng: data.longitude,
        accuracy: data.accuracy,
        charging: false,
        battery: data.battery
      })
    } catch(e) {
      this.setState({ error: e.message })
    }
  }

  componentDidMount() {
    this.updateData()

    let interval = setInterval(this.updateData.bind(this), 60000)
    this.setState({ interval })
  }

  componentWillUnmount() {
    if(this.state.interval) clearInterval(this.state.interval)
  }

  render() {
    if(this.state.error) {
      return <Error message={this.state.error} />
    }

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
          <a href={`https://www.google.com/maps/search/${this.state.lat},${this.state.lng}`} target="_blank">
            {this.state.location}
          </a>.
        </BigText>
        My battery was at
        <BigText>
          {this.state.battery}%
          {this.state.charging && <FontAwesomeIcon className="charging-icon" icon={faBolt} color="#0488f4" title="Charging"/>}
        </BigText>
        when this was reported.

        <Map center={[this.state.lat, this.state.lng]} zoom={17}>
          <TileLayer
            attribution='&amp;copy <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker position={[this.state.lat, this.state.lng]}/>
          <Circle center={[this.state.lat, this.state.lng]} radius={this.state.accuracy}/>
        </Map>

        <strong className="warning">Please note:</strong> Do not share this URL with anyone unless told otherwise.
      </div>
    )
  }
}

export default KeyPage
