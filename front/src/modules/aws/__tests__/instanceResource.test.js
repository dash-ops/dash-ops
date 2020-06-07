import http from "../../../helpers/http"
import { getInstances, startInstance, stopInstance } from "../instanceResource"

it("should return instances list", async () => {
  const mockResponse = [
    {
      instance_id: "666",
      name: "app-ops",
      status: "stopped",
    },
  ]
  jest.spyOn(http, "get").mockResolvedValue({
    data: mockResponse,
  })

  const resp = await getInstances()

  expect(http.get).toBeCalledWith(`${process.env.REACT_APP_API_URL}/v1/ec2/instances`)
  expect(resp.data).toEqual(mockResponse)
})

it("should return intance status when startInstance called", async () => {
  const mockResponse = { status: "running" }
  jest.spyOn(http, "post").mockResolvedValue({
    data: mockResponse,
  })

  const instanceID = 666
  const resp = await startInstance(instanceID)

  expect(http.post).toBeCalledWith(
    `${process.env.REACT_APP_API_URL}/v1/ec2/instance/start/${instanceID}`,
  )
  expect(resp.data).toEqual(mockResponse)
})

it("should return intance status when stopInstance called", async () => {
  const mockResponse = { status: "stopped" }
  jest.spyOn(http, "post").mockResolvedValue({
    data: mockResponse,
  })

  const instanceID = 666
  const resp = await stopInstance(instanceID)

  expect(http.post).toBeCalledWith(
    `${process.env.REACT_APP_API_URL}/v1/ec2/instance/stop/${instanceID}`,
  )
  expect(resp.data).toEqual(mockResponse)
})
