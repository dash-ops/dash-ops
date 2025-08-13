import { render, screen, cleanup, act } from "@testing-library/react"
import { MemoryRouter } from "react-router"
import { notification } from "antd"
import * as userResource from "../../modules/oauth2/userResource"

import Toolbar from "../Toolbar"

vi.mock("../../modules/oauth2/userResource")

const mockUser = { name: "Bla", avatar_url: "/image.jpg" }

afterEach(cleanup)

it("should show error notification when GitHub session is invalid", async () => {
  userResource.getUserData.mockRejectedValue(new Error())
  const notificationSpy = vi.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(
      <MemoryRouter>
        <Toolbar oAuth2 />
      </MemoryRouter>
    )
  })

  expect(notificationSpy).toHaveBeenCalledWith({ message: "Failed to fetch user data" })
  userResource.getUserData.mockRestore()
  notificationSpy.mockRestore()
})

it("should display user data when user is logged in", async () => {
  userResource.getUserData.mockResolvedValue({ data: mockUser })

  await act(async () => {
    render(
      <MemoryRouter>
        <Toolbar oAuth2 />
      </MemoryRouter>
    )
  })

  const avatarImg = document.querySelector('img[src="/image.jpg"]')
  expect(avatarImg.src).toContain(mockUser.avatar_url)
  expect(screen.getByText(mockUser.name)).toBeInTheDocument()
  userResource.getUserData.mockRestore()
})
