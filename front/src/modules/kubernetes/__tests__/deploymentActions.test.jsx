import { render, screen, cleanup } from "@testing-library/react"
import { MemoryRouter } from "react-router"
import userEvent from "@testing-library/user-event"
import DeploymentActions from "../DeploymentActions"

const mockDeployment = {
  name: "my-microservice",
  namespace: "default",
  pod_count: 0,
  pod_info: {
    current: 0,
  }
}

afterEach(cleanup)

it("should call deployment scale up when clicked on play button", async () => {
  const mockCp = { ...mockDeployment }
  const toUp = vi.fn()

  render(
    <MemoryRouter>
      <DeploymentActions deployment={mockCp} toUp={toUp} toDown={() => {}} />
    </MemoryRouter>
  )

  const upButton = screen.getByRole("button", { name: /Up/i })
  await userEvent.click(upButton)

  expect(toUp).toBeCalled()
})

it("should call deployment scale down when clicked on stop button", async () => {
  const mockCp = { ...mockDeployment, pod_count: 1, pod_info: { current: 1 } }
  const toDown = vi.fn()

  render(
    <MemoryRouter>
      <DeploymentActions deployment={mockCp} toUp={() => {}} toDown={toDown} />
    </MemoryRouter>
  )

  const downButton = screen.getByRole("button", { name: /Down/i })
  await userEvent.click(downButton)

  expect(toDown).toBeCalled()
})
