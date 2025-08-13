import { render, screen, cleanup, act } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { notification } from "antd"
import * as instanceResource from "../instanceResource"
import InstancePage from "../InstancePage"

vi.mock("axios", () => ({
  default: {
    create: () => ({
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      },
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn()
    })
  }
}))

afterEach(cleanup)

it("should display empty table when no instances are returned", async () => {
  vi.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: [] })

  await act(async () => {
    render(<InstancePage />)
  })


  await screen.findByRole("searchbox")
  expect(screen.getByRole("searchbox")).toBeInTheDocument()
  expect(screen.getByRole("button", { name: /clear/i })).toBeInTheDocument()
})

it("should display instances table with data when instances are available", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      status: "stopped",
    },
  ]
  vi.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />)
  })

  const tds = screen.getAllByRole("cell")
  expect(tds[0].textContent).toBe("app-ops")
  expect(tds[1].textContent).toBe("666")
})

it("should show error notification when instance fetch fails", async () => {
  vi.spyOn(instanceResource, "getInstances").mockRejectedValue(new Error())
  vi.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(<InstancePage />)
  })

  expect(notification.error).toBeCalledWith({
    message: "Ops... Failed to fetch API data",
  })
})

it("should filter instances when typing in search field", async () => {
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
  vi.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />)
  })

  const input = screen.getByRole("searchbox")
  await act(async () => {
    await userEvent.type(input, "app-42")
  })

  await screen.findByText("app-42")
  
  const tds = screen.getAllByRole("cell")
  expect(tds[0].textContent).toBe("app-42")
})

it("should display start button for stopped instances", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      state: "stopped",
    },
  ]
  
  vi.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />)
  })

  const startButton = screen.getByRole("button", { name: /start/i })
  expect(startButton).toBeInTheDocument()
  expect(startButton).not.toBeDisabled()
})

it("should display stop button for running instances", async () => {
  const mockInstances = [
    {
      instance_id: "666",
      name: "app-ops",
      state: "running",
    },
  ]
  
  vi.spyOn(instanceResource, "getInstances").mockResolvedValue({ data: mockInstances })

  await act(async () => {
    render(<InstancePage />)
  })

  const stopButton = screen.getByRole("button", { name: /stop/i })
  expect(stopButton).toBeInTheDocument()
  expect(stopButton).not.toBeDisabled()
})
