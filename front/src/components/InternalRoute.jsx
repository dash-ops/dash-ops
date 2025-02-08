import { useEffect } from "react"
import PropTypes from "prop-types"
import { Route, useNavigate } from "react-router"
import { verifyToken } from "../helpers/oauth"

function InternalRoute({ children, redirect, oAuth2, ...rest }) {
  const navigate = useNavigate()

  useEffect(() => {
    if (oAuth2 && !verifyToken()) {
      navigate(redirect, { state: { from: rest.location.pathname } })
    }
  }, [oAuth2, navigate, redirect, rest.location.pathname])

  return (
    <Route path={rest.path} exact={rest.exact}>
      {children}
    </Route>
  )
}

InternalRoute.propTypes = {
  children: PropTypes.element.isRequired,
  oAuth2: PropTypes.bool.isRequired,
  redirect: PropTypes.string,
}

InternalRoute.defaultProps = {
  redirect: "/login",
}

export default InternalRoute
