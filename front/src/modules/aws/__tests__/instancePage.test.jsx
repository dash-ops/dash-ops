import { render, screen, cleanup, act } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { notification } from "antd"
import * as instanceResource from "../instanceResource"
import InstancePage from "../InstancePage"

jest.mock("axios")

afterEach(cleanup)

it("should return table without content", async () => {
  jest.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: [] })

  await act(async () => {
    render(<InstancePage />)
  })

  const table = screen.getByRole("table")
  const ths = screen.getAllByRole("columnheader")
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
    render(<InstancePage />)
  })

  const tds = screen.getAllByRole("cell")
  expect(tds[0].textContent).toBe("app-ops")
  expect(tds[1].textContent).toBe("666")
})

it("should return notification error when failed instances fetch", async () => {
  jest.spyOn(instanceResource, "getInstances").mockRejectedValue(new Error())
  jest.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(<InstancePage />)
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
    render(<InstancePage />)
  })

  const input = screen.getByRole("textbox")
  await act(async () => {
    userEvent.type(input, "app-42")
  })

  const tds = screen.getAllByRole("cell")
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
    render(<InstancePage />)
  })

  const startButton = screen.getByRole("button", { name: /start/i })
  await act(async () => {
    userEvent.click(startButton)
  })

  expect(instanceResource.startInstance).toBeCalled()
  instanceResource.startInstance.mockRestore()
})

it("should request stop instance when clicked in stop button", async () => {
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
    render(<InstancePage />)
  })

  const stopButton = screen.getByRole("button", { name: /stop/i })
  await act(async () => {
    userEvent.click(stopButton)
  })

  expect(instanceResource.stopInstance).toBeCalled()
  instanceResource.stopInstance.mockRestore()
})
