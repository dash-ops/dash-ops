import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import InstanceTag from "../InstanceTag"

let container = null
beforeEach(() => {
  container = document.createElement("div")
  document.body.appendChild(container)
})

afterEach(() => {
  unmountComponentAtNode(container)
  container.remove()
  container = null
})

it("should return tag when instance state passed", () => {
  act(() => {
    render(<InstanceTag state="running" />, container)
  })

  const tag = container.querySelector("span")
  expect(tag.textContent).toBe("running")
  expect(tag.className).toBe("ant-tag ant-tag-green")
})
