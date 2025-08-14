import { useEffect, useState, useCallback, useRef } from "react"
import PropTypes from "prop-types"
import { useNavigate, Link } from "react-router"
import { Menu, Avatar, notification } from "antd"
import { LogoutOutlined, UserOutlined } from "@ant-design/icons"
import { cleanToken } from "../helpers/oauth"
import { getUserData } from "../modules/oauth2/userResource"

function Toolbar({ oAuth2 }) {
  const [user, setUser] = useState()
  const navigate = useNavigate()
  const userDataFetched = useRef(false)

  const logout = useCallback(() => {
    setUser(null)
    cleanToken()
    userDataFetched.current = false // Reset cache on logout
    navigate("/login")
  }, [navigate])

  useEffect(() => {
    if (!oAuth2 || userDataFetched.current) {
      return
    }

    async function fetchData() {
      try {
        const result = await getUserData()
        setUser(result.data)
        userDataFetched.current = true // Mark as fetched
      } catch {
        notification.error({ message: "Failed to fetch user data" })
      }
    }

    fetchData()
  }, [oAuth2])

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
                key: "profile",
                label: (
                  <Link to="/profile" style={{ color: 'inherit', textDecoration: 'none' }}>
                    <UserOutlined />
                    Profile
                  </Link>
                ),
              },
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
