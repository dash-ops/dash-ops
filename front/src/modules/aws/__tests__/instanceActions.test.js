import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import { notification } from "antd"
import InstanceActions from "../InstanceActions"

let container = null
const mockInstance = {
  instance_id: "666",
  name: "xpto",
  platform: "",
}
beforeEach(() => {
  container = document.createElement("div")
  document.body.appendChild(container)
})

afterEach(() => {
  unmountComponentAtNode(container)
  container.remove()
  container = null
})

it("should call start instance when clicked on play button", () => {
  const mockCp = { ...mockInstance, state: "stopped" }
  const toStart = jest.fn()

  act(() => {
    render(<InstanceActions instance={mockCp} toStart={toStart} toStop={() => {}} />, container)
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[1].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(toStart).toBeCalled()
})

it("should call stop instance when clicked on stop button", () => {
  const mockCp = { ...mockInstance, state: "running" }
  const toStop = jest.fn()

  act(() => {
    render(<InstanceActions instance={mockCp} toStart={() => {}} toStop={toStop} />, container)
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[1].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(toStop).toBeCalled()
})

it("should return new location address when clicked on ssh button", () => {
  const originalLocation = window.location
  delete window.location
  window.location = new URL("http://mock-location.com")

  act(() => {
    render(
      <InstanceActions instance={mockInstance} toStart={() => {}} toStop={() => {}} />,
      container,
    )
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[0].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(window.location).toBe("ssh://xpto")
  window.location = originalLocation
})

it("should return notification when clicked on ssh button and instance plataform is windows", () => {
  mockInstance.platform = "windows"
  jest.spyOn(notification, "error").mockImplementation(() => {})

  act(() => {
    render(
      <InstanceActions instance={mockInstance} toStart={() => {}} toStop={() => {}} />,
      container,
    )
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[0].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(notification.error).toBeCalledWith({
    message: `Sorry... I'm afraid I can't do that...`,
    description: `
        Windows does not provides a method to connect a Remote Desktop via URL.
        You can try to connect via command line using on Windows: mstsc /v:${mockInstance.name}
      `,
  })
  mockInstance.platform = ""
})
