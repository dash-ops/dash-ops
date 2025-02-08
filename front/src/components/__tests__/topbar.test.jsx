import { render, screen, cleanup, act } from "@testing-library/react"
import { MemoryRouter } from "react-router"
import { notification } from "antd"
import * as userResource from "../../modules/oauth2/userResource"

import Toolbar from "../Toolbar"

vi.mock("../../modules/oauth2/userResource")

const mockUser = { name: "Bla", avatar_url: "/image.jpg" }

afterEach(cleanup)

it("should notify error when session on github invalidates", async () => {
  userResource.getUserData.mockRejectedValue(new Error())

  await act(async () => {
    render(
      <MemoryRouter>
        <Toolbar oAuth2 />
      </MemoryRouter>
    )
  })

  expect(notification.error).toBeCalledWith({ message: "Failed to fetch user data" })
  userResource.getUserData.mockRestore()
})

it("should user data when logged in user", async () => {
  userResource.getUserData.mockResolvedValue({ data: mockUser })

  await act(async () => {
    render(
      <MemoryRouter>
        <Toolbar oAuth2 />
      </MemoryRouter>
    )
  })

  expect(screen.getByRole("img").src).toBe(mockUser.avatar_url)
  expect(screen.getByText(mockUser.name)).toBeInTheDocument()
  userResource.getUserData.mockRestore()
})
