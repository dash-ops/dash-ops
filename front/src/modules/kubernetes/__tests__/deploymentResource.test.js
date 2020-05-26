import http from "../../../helpers/http"
import { getDeployments, upDeployment, downDeployment } from "../deploymentResource"

it("should return deployments list", async () => {
  const mockResponse = [
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
  jest.spyOn(http, "get").mockResolvedValue({
    data: mockResponse,
  })

  const resp = await getDeployments()

  expect(http.get).toBeCalledWith(`${process.env.REACT_APP_API_URL}/v1/k8s/deployments`)
  expect(resp.data).toEqual(mockResponse)
})

it("should start pod when upDeployment called", async () => {
  const mockResponse = {}
  jest.spyOn(http, "post").mockResolvedValue({
    data: mockResponse,
  })

  const name = "my-microservice"
  const namespace = "default"
  const resp = await upDeployment(name, namespace)

  expect(http.post).toBeCalledWith(
    `${process.env.REACT_APP_API_URL}/v1/k8s/deployment/up/${namespace}/${name}`,
  )
  expect(resp.data).toEqual(mockResponse)
})

it("should stop pod when downDeployment called", async () => {
  const mockResponse = {}
  jest.spyOn(http, "post").mockResolvedValue({
    data: mockResponse,
  })

  const name = "my-microservice"
  const namespace = "default"
  const resp = await downDeployment(name, namespace)

  expect(http.post).toBeCalledWith(
    `${process.env.REACT_APP_API_URL}/v1/k8s/deployment/down/${namespace}/${name}`,
  )
  expect(resp.data).toEqual(mockResponse)
})
