import { useEffect, useState, useCallback } from "react"
import PropTypes from "prop-types"
import { useNavigate } from "react-router"
import { Menu, Avatar, notification } from "antd"
import { LogoutOutlined } from "@ant-design/icons"
import { cleanToken } from "../helpers/oauth"
import { getUserData } from "../modules/oauth2/userResource"

function Toolbar({ oAuth2 }) {
  const [user, setUser] = useState()
  const navigate = useNavigate()

  const logout = useCallback(() => {
    setUser(null)
    cleanToken()
    navigate("/login")
  }, [navigate])

  useEffect(() => {
    if (!oAuth2) {
      return
    }

    async function fetchData() {
      try {
        const result = await getUserData()
        setUser(result.data)
      } catch {
        notification.error({ message: "Failed to fetch user data" })
      }
    }

    fetchData()
  }, [logout, oAuth2])

  const menuItems = user
    ? [
        {
          key: "user",
          label: (
            <>
              <Avatar src={user.avatar_url} />
              <span style={{ padding: "10px" }}>{user.name}</span>
            </>
          ),
          children: [
            {
              key: "logout",
              label: (
                <>
                  <LogoutOutlined />
                  Logout
                </>
              ),
              onClick: logout,
            },
          ],
        },
      ]
    : []

  return (
    <Menu mode="horizontal" style={{ lineHeight: "64px", float: "right" }} items={menuItems} />
  )
}

Toolbar.propTypes = {
  oAuth2: PropTypes.bool.isRequired,
}

export default Toolbar
