import React from "react"
import { GithubOutlined } from "@ant-design/icons"

function Footer() {
  return (
    <>
      <a href="https://github.com/dash-ops/dash-ops" alt="DashOps Repository">
        <GithubOutlined />
      </a>{" "}
      <a href="https://dash-ops.github.io/" alt="DashOps WebSite">
        DashOPS
      </a>{" "}
      ©{new Date().getFullYear()}
    </>
  )
}

export default Footer
