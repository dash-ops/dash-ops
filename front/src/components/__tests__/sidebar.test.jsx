import { render, screen, cleanup, act } from "@testing-library/react"
import { MemoryRouter } from "react-router"
import userEvent from "@testing-library/user-event"

import Sidebar from "../Sidebar"

const mockMenus = [
  {
    label: "AWS Plugin",
    key: "aws"
  },
  {
    label: "Kubernetes Plugin",
    key: "k8s"
  },
]

afterEach(cleanup)

it("should render sidebar with menu items", () => {
  act(() => {
    render(
      <MemoryRouter>
        <Sidebar menus={mockMenus} />
      </MemoryRouter>
    )
  })

  expect(screen.getByText("AWS Plugin")).toBeInTheDocument()
  expect(screen.getByText("Kubernetes Plugin")).toBeInTheDocument()
})

it("should highlight menu item when clicked", async () => {
  act(() => {
    render(
      <MemoryRouter>
        <Sidebar menus={mockMenus} />
      </MemoryRouter>
    )
  })
  const firstItem = screen.getByText("AWS Plugin").parentElement
  await userEvent.click(firstItem)

  expect(firstItem).toHaveClass("ant-menu-item-active")
})
