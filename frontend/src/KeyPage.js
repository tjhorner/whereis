import React from 'react'
import './App.css'
import BigText from './BigText'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBolt } from '@fortawesome/free-solid-svg-icons'
import axios from "axios"
import moment from "moment"
import Error from './Error'
import { Map, TileLayer, Marker, Circle, Polyline } from "react-leaflet"

class KeyPage extends React.Component {
  state = {
    loading: true,
    locations: [ ],
    currentLocation: 0,
    interval: null,
    error: null
  }

  get location() {
    if(this.state.locations.length === 0) return

    let loc = this.state.locations[this.state.currentLocation]
    let at = moment(loc.at)

    return {
      lastSeen: at.fromNow(),
      lastSeenAbsolute: at.format("dddd, MMMM Do YYYY, h:mm:ss a"),
      location: loc.coarse_location,
      lat: loc.latitude,
      lng: loc.longitude,
      accuracy: loc.accuracy,
      charging: false,
      battery: loc.battery
    }
  }

  get latLngs() {
    return this.state.locations.map(l => [l.latitude, l.longitude])
  }

  async updateData() {
    try {
      const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/v1/locations`, {
        params: { key: this.props.match.params.id }
      })

      if(data.error) {
        this.setState({ error: data.error })
        return
      }

      this.setState({
        loading: false,
        locations: data
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

  updateSlider(ev) {
    this.setState({ currentLocation: parseInt(ev.target.value) })
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

    let points
    if(this.state.locations.length > 1) {
      points = this.state.locations.map((loc, currentLocation) => {
        return (
          <Circle
            key={currentLocation}
            color="#388e3c"
            onClick={() => this.setState({ currentLocation })}
            center={[loc.latitude, loc.longitude]}
            fillOpacity={1}
            radius={3}/>
        )
      })
    }

    return (
      <div className="App">
        According to my phone, I was {this.state.currentLocation === 0 && "last"} seen
        <BigText title={this.location.lastSeenAbsolute}>{this.location.lastSeen}</BigText>
        somewhere in
        <BigText>
          <a href={`https://www.google.com/maps/search/${this.location.lat},${this.location.lng}`} target="_blank">
            {this.location.location}
          </a>.
        </BigText>
        My battery was at
        <BigText>
          {this.location.battery}%
          {this.location.charging && <FontAwesomeIcon className="charging-icon" icon={faBolt} color="#0488f4" title="Charging"/>}
        </BigText>
        when this was reported.

        {this.state.locations.length > 1 && <input className="slider" type="range" min="0" max={this.state.locations.length - 1} value={this.state.currentLocation} onChange={this.updateSlider.bind(this)}/>}

        <Map center={[this.location.lat, this.location.lng]} zoom={17}>
          <TileLayer
            attribution='&amp;copy <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker position={[this.location.lat, this.location.lng]}/>
          <Circle center={[this.location.lat, this.location.lng]} radius={this.location.accuracy}/>
          <Polyline positions={this.latLngs}/>
          {points}
        </Map>

        <strong className="warning">Please note:</strong> Do not share this URL with anyone unless told otherwise.
      </div>
    )
  }
}

export default KeyPage
