import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import { notification } from "antd"
import * as instanceResource from "../instanceResource"
import InstancePage from "../InstancePage"

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

it("should return table without content", async () => {
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: [] })

  await act(async () => {
    render(<InstancePage />, container)
  })

  const table = container.querySelector("table")
  const ths = table.querySelectorAll("th")
  expect(ths[0].textContent).toBe("Instance")
  expect(ths[1].textContent).toBe("Instance Id")
  expect(ths[2].textContent).toBe("State")
  expect(ths[3].textContent).toBe("Action")
  const tbody = table.querySelector("tbody")
  expect(tbody.textContent).toBe("")
})

it("should return table with content", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      status: "stopped",
    },
  ]
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />, container)
  })

  const tbody = container.querySelector("tbody")
  const tds = tbody.querySelectorAll("td")
  expect(tds[0].textContent).toBe("app-ops")
  expect(tds[1].textContent).toBe("666")
})

it("should return notification error when failed instances fetch", async () => {
  jest.spyOn(instanceResource, "getInstances").mockRejectedValue(new Error())
  jest.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(<InstancePage />, container)
  })

  expect(notification.error).toBeCalledWith({
    message: "Ops... Failed to fetch API data",
  })
})

it("should filter the list when filled in the search field", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      state: "stopped",
    },
    {
      instance_id: "999",
      name: "app-42",
      state: "running",
    },
  ]
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />, container)
  })

  const input = container.querySelector("input")
  act(() => {
    Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value").set.call(input, "app-42")
    input.dispatchEvent(new Event("input", { bubbles: true }))
  })

  const tbody = container.querySelector("tbody")
  const tds = tbody.querySelectorAll("td")
  expect(tds[0].textContent).toBe("app-42")
})

it("should request start instance when clicked in start button", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      state: "stopped",
    },
  ]
  const mockStartResponse = {
    current_state: "pending",
  }
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })
  jest.spyOn(instanceResource, "startInstance").mockResolvedValue({ data: mockStartResponse })

  await act(async () => {
    render(<InstancePage />, container)
  })

  const tbody = container.querySelector("tbody")
  const actionTd = tbody.querySelectorAll("td")[3]
  const startButton = actionTd.querySelectorAll("button")[1]
  await act(async () => {
    startButton.dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(instanceResource.startInstance).toBeCalled()
  instanceResource.startInstance.mockRestore()
})

it("should request start instance when clicked in stop button", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      state: "running",
    },
  ]
  const mockStopResponse = {
    current_state: "stopping",
  }
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })
  jest.spyOn(instanceResource, "stopInstance").mockResolvedValue({ data: mockStopResponse })

  await act(async () => {
    render(<InstancePage />, container)
  })

  const tbody = container.querySelector("tbody")
  const actionTd = tbody.querySelectorAll("td")[3]
  const stopButton = actionTd.querySelectorAll("button")[1]
  await act(async () => {
    stopButton.dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(instanceResource.stopInstance).toBeCalled()
  instanceResource.stopInstance.mockRestore()
})
