import { render, unmountComponentAtNode } from "react-dom"
import { act } from "react-dom/test-utils"
import { notification } from "antd"
import * as deploymentResource from "../deploymentResource"
import * as namespaceResource from "../namespaceResource"
import DeploymentPage from "../DeploymentPage"

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

it("should return table without content", async () => {
  jest.spyOn(deploymentResource, "getDeployments").mockResolvedValue({ data: [] })

  await act(async () => {
    render(<DeploymentPage />, container)
  })

  const table = container.querySelector("table")
  const ths = table.querySelectorAll("th")
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
    render(<DeploymentPage />, container)
  })

  const tbody = container.querySelector("tbody")
  const tds = tbody.querySelectorAll("td")
  expect(tds[0].textContent).toBe("my-microservice")
  expect(tds[1].textContent).toBe("default")
})

it("should return notification error when failed instances fetch", async () => {
  jest.spyOn(namespaceResource, "getNamespaces").mockResolvedValue({ data: [{ name: "default" }] })
  jest.spyOn(deploymentResource, "getDeployments").mockRejectedValue(new Error())
  jest.spyOn(notification, "error").mockImplementation(() => {})

  await act(async () => {
    render(<DeploymentPage />, container)
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
    render(<DeploymentPage />, container)
  })

  const input = container.querySelector("input")
  act(() => {
    Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value").set.call(
      input,
      "accountancy",
    )
    input.dispatchEvent(new Event("input", { bubbles: true }))
  })

  const tbody = container.querySelector("tbody")
  const tds = tbody.querySelectorAll("td")
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
    render(<DeploymentPage />, container)
  })

  const tbody = container.querySelector("tbody")
  const actionTd = tbody.querySelectorAll("td")[3]
  const upButton = actionTd.querySelector("button")
  await act(async () => {
    upButton.dispatchEvent(new MouseEvent("click", { bubbles: true }))
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
    render(<DeploymentPage />, container)
  })

  const tbody = container.querySelector("tbody")
  const actionTd = tbody.querySelectorAll("td")[3]
  const downButton = actionTd.querySelector("button")
  await act(async () => {
    downButton.dispatchEvent(new MouseEvent("click", { bubbles: true }))
  })

  expect(deploymentResource.downDeployment).toBeCalled()
  deploymentResource.downDeployment.mockRestore()
})
