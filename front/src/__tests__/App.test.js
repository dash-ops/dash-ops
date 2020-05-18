import React from "react"
import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import * as userResource from "../modules/oauth2/userResource"
import * as oauth from "../helpers/oauth"
import App from "../App"

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

it.skip("should render instances page when logged in user", async () => {
  jest.spyOn(oauth, "verifyToken").mockReturnValue(true)
  jest.spyOn(userResource, "getUserData").mockResolvedValue({ name: "Bla" })

  await act(async () => {
    render(<App />, container)
  })

  // const imgLogo = container.querySelector(".logo").childNodes[0]
  // expect(imgLogo.alt).toBe("DashOps")
  const footer = container.querySelector(".footer")
  expect(footer.textContent).toBe(`DashOPS ©${new Date().getFullYear()}`)
})
