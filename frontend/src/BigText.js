import React from 'react';
import './BigText.css';

const BigText = props => (
  <div className="big-text" title={props.title}>{props.children}</div>
)

export default BigText
