import { render, screen, cleanup, act } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { notification } from "antd"
import * as deploymentResource from "../deploymentResource"
import * as namespaceResource from "../namespaceResource"
import DeploymentPage from "../DeploymentPage"

jest.mock("axios")

afterEach(cleanup)

it("should return table without content", async () => {
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: [] })

  await act(async () => {
    render(<DeploymentPage />)
  })

  const table = screen.getByRole("table")
  const ths = screen.getAllByRole("columnheader")
  expect(ths[0].textContent).toBe("Name")
  expect(ths[1].textContent).toBe("Namespace")
  expect(ths[2].textContent).toBe("Pods running")
  expect(ths[3].textContent).toBe("Action")
  const tbody = table.querySelector("tbody")
  expect(tbody.textContent).toBe("")
})

it("should return table with content", async () => {
  const mockDeployments = [
    {
      name: "my-microservice",
      namespace: "default",
      pod_count: 0,
    },
  ]
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: mockDeployments })

  await act(async () => {
    render(<DeploymentPage />)
  })

  const tds = screen.getAllByRole("cell")
  expect(tds[0].textContent).toBe("my-microservice")
  expect(tds[1].textContent).toBe("default")
})

it("should return notification error when failed instances fetch", async () => {
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockRejectedValue(new Error())
  jest.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(<DeploymentPage />)
  })

  expect(notification.error).toBeCalledWith({
    message: "Ops... Failed to fetch API data",
  })
})

it("should filter the list when filled in the search field", async () => {
  const mockDeployments = [
    {
      name: "my-microservice",
      namespace: "default",
      pod_count: 0,
    },
    {
      name: "other-microservice",
      namespace: "default",
      pod_count: 0,
    },
  ]
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: mockDeployments })

  await act(async () => {
    render(<DeploymentPage />)
  })

  const input = screen.getByRole("textbox")
  await act(async () => {
    userEvent.type(input, "accountancy")
  })

  const tds = screen.getAllByRole("cell")
  expect(tds[0].textContent).toBe("other-microservice")
})

it("should request scale up deployment when clicked in up button", async () => {
  const mockDeployments = [
    {
      name: "my-microservice",
      namespace: "default",
      pod_count: 0,
    },
  ]
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: mockDeployments })
  jest.spyOn(deploymentResource, "upDeployment")

  await act(async () => {
    render(<DeploymentPage />)
  })

  const upButton = screen.getByRole("button", { name: /Up/i })
  await act(async () => {
    userEvent.click(upButton)
  })

  expect(deploymentResource.upDeployment).toBeCalled()
  deploymentResource.upDeployment.mockRestore()
})

it("should request scale down deployment when clicked in down button", async () => {
  const mockDeployments = [
    {
      name: "my-microservice",
      namespace: "default",
      pod_count: 1,
    },
  ]
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: mockDeployments })
  jest.spyOn(deploymentResource, "downDeployment")

  await act(async () => {
    render(<DeploymentPage />)
  })

  const downButton = screen.getByRole("button", { name: /Down/i })
  await act(async () => {
    userEvent.click(downButton)
  })

  expect(deploymentResource.downDeployment).toBeCalled()
  deploymentResource.downDeployment.mockRestore()
})
