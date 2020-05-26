import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import DeploymentActions from "../DeploymentActions"

let container = null
const mockDeployment = {
  name: "my-microservice",
  namespace: "default",
  pod_count: 0,
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

it("should call deployment scale up when clicked on play button", () => {
  const mockCp = { ...mockDeployment }
  const toUp = jest.fn()

  act(() => {
    render(<DeploymentActions deployment={mockCp} toUp={toUp} toDown={() => {}} />, container)
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[0].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(toUp).toBeCalled()
})

it("should call deployment scale down when clicked on stop button", () => {
  const mockCp = { ...mockDeployment, pod_count: 1 }
  const toDown = jest.fn()

  act(() => {
    render(<DeploymentActions deployment={mockCp} toUp={() => {}} toDown={toDown} />, container)
  })

  const buttons = document.querySelectorAll("button")
  act(() => {
    buttons[0].dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(toDown).toBeCalled()
})
