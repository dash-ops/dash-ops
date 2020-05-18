import React from "react"
import PropTypes from "prop-types"
import { Route, Redirect } from "react-router-dom"
import { verifyToken } from "../helpers/oauth"

function PrivateRoute({ children, redirect, ...rest }) {
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

PrivateRoute.propTypes = {
  children: PropTypes.element.isRequired,
  redirect: PropTypes.string,
}

PrivateRoute.defaultProps = {
  redirect: "/login",
}

export default PrivateRoute
