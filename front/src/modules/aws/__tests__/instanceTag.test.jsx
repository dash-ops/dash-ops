import { render, screen, cleanup } from "@testing-library/react"
import InstanceTag from "../InstanceTag"

afterEach(cleanup)

it("should return tag when instance state passed", () => {
  render(<InstanceTag state="running" />)

  const tag = screen.getByText("running")
  expect(tag).toBeInTheDocument()
  expect(tag).toHaveClass("ant-tag ant-tag-green")
})
