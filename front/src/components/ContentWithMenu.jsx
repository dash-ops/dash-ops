import { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { Routes, Route, useParams, useLocation, useNavigate } from "react-router"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

function ContentWithMenu({ pages, paramName = "context" }) {
  const params = useParams()
  const location = useLocation()
  const navigate = useNavigate()
  const [current, setCurrent] = useState(location.pathname)

  // Get the parameter value based on paramName
  const paramValue = params[paramName]

  useEffect(() => {
    setCurrent(location.pathname)
  }, [location.pathname])

  const onClick = (path) => {
    navigate(path)
  }

  const menuItems = pages
    .filter((page) => page.menu)
    .map((page) => ({
      path: page.path.replace(`:${paramName}`, paramValue),
      name: page.name,
    }))

  return (
    <div className="grid grid-cols-1 md:grid-cols-5 lg:grid-cols-6 xl:grid-cols-7 gap-4 px-2 py-4">
      <div className="md:col-span-1 xl:col-span-1">
        <nav className="space-y-1">
          {menuItems.map((item) => (
            <Button
              key={item.path}
              variant={current === item.path ? "default" : "ghost"}
              className={cn(
                "w-full justify-start",
                current === item.path && "bg-primary text-primary-foreground"
              )}
              onClick={() => onClick(item.path)}
            >
              {item.name}
            </Button>
          ))}
        </nav>
      </div>
      <div className="md:col-span-4 lg:col-span-5 xl:col-span-6">
        <Routes>
          {pages.map((page) => {
            const path = page.path.replace(`:${paramName}`, paramValue)
            const route = page.path.split(`:${paramName}`).pop()
            return (
              <Route
                key={path}
                path={route}
                element={page.element}
              />
            )
          })}
        </Routes>
      </div>
    </div>
  )
}

ContentWithMenu.propTypes = {
  pages: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired,
      menu: PropTypes.bool,
      element: PropTypes.object.isRequired,
    })
  ).isRequired,
  paramName: PropTypes.string, // Nome do parâmetro da URL (ex: "context", "key")
}

export default ContentWithMenu
