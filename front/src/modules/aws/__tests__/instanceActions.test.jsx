import { render, screen, cleanup } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { notification } from "antd"
import InstanceActions from "../InstanceActions"

const mockInstance = {
  instance_id: "666",
  name: "xpto",
  platform: "",
}

afterEach(cleanup)

it("should call start instance when clicked on play button", async () => {
  const mockCp = { ...mockInstance, state: "stopped" }
  const toStart = jest.fn()

  render(<InstanceActions instance={mockCp} toStart={toStart} toStop={() => {}} />)

  const playButton = screen.getByRole("button", { name: /Start/i })
  await userEvent.click(playButton)

  expect(toStart).toBeCalled()
})

it("should call stop instance when clicked on stop button", async () => {
  const mockCp = { ...mockInstance, state: "running" }
  const toStop = jest.fn()

  render(<InstanceActions instance={mockCp} toStart={() => {}} toStop={toStop} />)

  const stopButton = screen.getByRole("button", { name: /Stop/i })
  await userEvent.click(stopButton)

  expect(toStop).toBeCalled()
})

it("should return new location address when clicked on ssh button", async () => {
  const originalLocation = window.location
  delete window.location
  window.location = new URL("http://mock-location.com")

  render(<InstanceActions instance={mockInstance} toStart={() => {}} toStop={() => {}} />)

  const sshButton = screen.getByRole("button", { name: /desktop/i })
  await userEvent.click(sshButton)

  expect(window.location).toBe("ssh://xpto")
  window.location = originalLocation
})

it("should return notification when clicked on ssh button and instance plataform is windows", async () => {
  const mockCp = { ...mockInstance, platform: "windows" }
  jest.spyOn(notification, "error").mockImplementation(() => {})

  render(<InstanceActions instance={mockCp} toStart={() => {}} toStop={() => {}} />)

  const sshButton = screen.getByRole("button", { name: /desktop/i })
  await userEvent.click(sshButton)

  expect(notification.error).toBeCalledWith({
    message: `Sorry... I'm afraid I can't do that...`,
    description: `
        Windows does not provides a method to connect a Remote Desktop via URL.
        You can try to connect via command line using on Windows: mstsc /v:${mockInstance.name}
      `,
  })
  mockInstance.platform = ""
})
