import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import useInterval from "../useInterval"

jest.useFakeTimers()

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

it("should callback after delay", () => {
  const callback = jest.fn()

  act(() => {
    render(
      <TestHookComponent
        callback={() => {
          useInterval(callback, 100)
        }}
      />,
      container,
    )
  })

  expect(callback).toHaveBeenCalledTimes(0)
  jest.advanceTimersByTime(400)
  expect(callback).toHaveBeenCalledTimes(4)
})

const TestHookComponent = ({ callback }) => {
  callback()
  return null
}
