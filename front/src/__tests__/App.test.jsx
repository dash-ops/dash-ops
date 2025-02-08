import { render, screen, cleanup, act } from "@testing-library/react"
import * as userResource from "../modules/oauth2/userResource"
import * as oauth from "../helpers/oauth"
import App from "../App"

vi.mock("../modules/oauth2/userResource")
vi.mock("../helpers/oauth")

afterEach(cleanup)

it.skip("should render instances page when logged in user", async () => {
  oauth.verifyToken.mockReturnValue(true)
  userResource.getUserData.mockResolvedValue({ name: "Bla" })

  await act(async () => {
    render(<App />)
  })

  const footer = screen.getByText(`DashOPS ©${new Date().getFullYear()}`)
  expect(footer).toBeInTheDocument()
})
