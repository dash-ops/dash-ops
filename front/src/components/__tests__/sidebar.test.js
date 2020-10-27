import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import { MemoryRouter } from "react-router-dom"

import Sidebar from "../Sidebar"

let container = null
const mockMenus = [
  {
    name: "AWS Plugin",
    path: "/aws/ec2",
    component: () => <div>Hello AWS!</div>,
  },
  {
    name: "Kubernetes Plugin",
    path: "/k8s/deployment",
    component: () => <div>Hello K8S!</div>,
  },
]

beforeEach(() => {
  container = document.createElement("div")
  document.body.appendChild(container)
})

afterEach(() => {
  unmountComponentAtNode(container)
  container.remove()
  container = null
})

it("renders without itens", () => {
  act(() => {
    render(
      <MemoryRouter>
        <Sidebar />
      </MemoryRouter>,
      container,
    )
  })

  expect(container.textContent).toBe("DashOPS")
})

it("renders with itens", () => {
  act(() => {
    render(
      <MemoryRouter>
        <Sidebar menus={mockMenus} />
      </MemoryRouter>,
      container,
    )
  })

  const links = container.querySelectorAll("a")
  expect(links[0].textContent).toBe("AWS Plugin")
  expect(links[1].textContent).toBe("Kubernetes Plugin")
})

it("should change current value when clicked on menu item", () => {
  act(() => {
    render(
      <MemoryRouter>
        <Sidebar menus={mockMenus} />
      </MemoryRouter>,
      container,
    )
  })

  const firstItem = container.querySelector("li")
  act(() => {
    firstItem.dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(firstItem.className.includes("ant-menu-item-selected")).toBeTruthy()
})
