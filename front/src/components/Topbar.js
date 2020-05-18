import React, { useEffect, useState, useCallback } from "react"
import { useHistory } from "react-router-dom"
import { Menu, Avatar, notification } from "antd"
import { LogoutOutlined } from "@ant-design/icons"
import { cleanToken } from "../helpers/oauth"
import { getUserData } from "../modules/oauth2/userResource"

function Topbar() {
  const [user, setUser] = useState()
  const history = useHistory()

  const logout = useCallback(() => {
    setUser(null)
    cleanToken()
    history.push("/login")
  }, [history])

  useEffect(() => {
    async function fetchData() {
      try {
        const result = await getUserData()
        setUser(result.data)
      } catch (e) {
        notification.error({ message: "Failed to fetch user data" })
      }
    }

    fetchData()
  }, [logout])

  return (
    <Menu mode="horizontal" style={{ lineHeight: "64px", float: "right" }}>
      {user && (
        <Menu.SubMenu
          title={
            <>
              <Avatar src={user.avatar_url} />
              <span style={{ padding: "10px" }}>{user.name}</span>
            </>
          }
        >
          <Menu.Item key="logout" onClick={logout}>
            <LogoutOutlined />
            Logout
          </Menu.Item>
        </Menu.SubMenu>
      )}
    </Menu>
  )
}

export default Topbar
