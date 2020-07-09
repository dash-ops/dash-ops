import React from "react"
import ReactDOM from "react-dom"
import { BrowserRouter as Router } from "react-router-dom"
import "antd/dist/antd.css"
import App from "./App"

ReactDOM.render(
  <Router>
    <App />
  </Router>,
  document.getElementById("root"),
)
