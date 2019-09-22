import React from 'react';
import './Error.css';
import hyena from "./hyena.png"
import { Link } from "react-router-dom"

const Error = props => (
  <div className="error">
    <figure>
      <img src={hyena}/>
      <figcaption>(Hyena sticker from <a href="https://store.line.me/stickershop/product/4439542/en" target="_blank" rel="noopener noreferrer">this pack</a>)</figcaption>
    </figure>

    <h1>Something went wrong.</h1>
    <h4>{props.message}</h4>
    {window.location.pathname !== "/" && <Link to="/">Go back</Link>}
  </div>
)

export default Error
