import React from 'react';
import './Error.css';
import hyena from "./hyena.png"
import { Link } from "react-router-dom"

const Error = props => (
  <div className="error">
    <img src={hyena}/>
    <h1>Something went wrong.</h1>
    <h4>{props.message}</h4>
    {window.location.pathname !== "/" && <Link to="/">Go back</Link>}
  </div>
)

export default Error
