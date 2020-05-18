import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import { MemoryRouter } from "react-router-dom"
import { notification } from "antd"
import * as userResource from "../../modules/oauth2/userResource"

import Topbar from "../Topbar"

let container = null
const mockUser = { name: "Bla", avatar_url: "/image.jpg" }

beforeEach(() => {
  container = document.createElement("div")
  document.body.appendChild(container)
})

afterEach(() => {
  unmountComponentAtNode(container)
  container.remove()
  container = null
})

it("should notify error when session on github invalidates", async () => {
  jest.spyOn(userResource, "getUserData").mockRejectedValue(new Error())
  jest.spyOn(notification, "error")

  await act(async () => {
    render(
      <MemoryRouter>
        <Topbar />
      </MemoryRouter>,
      container,
    )
  })

  expect(notification.error).toBeCalledWith({ message: "Failed to fetch user data" })
  userResource.getUserData.mockRestore()
})

it("should user data when logged in user", async () => {
  jest.spyOn(userResource, "getUserData").mockResolvedValue({ data: mockUser })

  await act(async () => {
    render(
      <MemoryRouter>
        <Topbar />
      </MemoryRouter>,
      container,
    )
  })

  const userItem = container.querySelectorAll("li")[1].querySelectorAll("span")
  expect(userItem[0].getElementsByTagName("img")[0].src).toBe(
    `http://localhost${mockUser.avatar_url}`,
  )
  expect(userItem[1].textContent).toBe(mockUser.name)
  userResource.getUserData.mockRestore()
})
