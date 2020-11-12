import PropTypes from "prop-types"
import { Route, Redirect } from "react-router-dom"
import { verifyToken } from "../helpers/oauth"

function InternalRoute({ children, redirect, oAuth2, ...rest }) {
  if (!oAuth2) {
    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Route {...rest}>{children}</Route>
  }

  return (
    <Route
      // eslint-disable-next-line react/jsx-props-no-spreading
      {...rest}
      render={({ location }) =>
        verifyToken() ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: redirect,
              state: location.pathname,
            }}
          />
        )
      }
    />
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
