import { useState, useEffect, useReducer, useCallback } from "react"
import { useParams } from "react-router-dom"
import { Row, Col, Table, Button, Input, notification, Form, Tag, Select } from "antd"
import { cancelToken } from "../../helpers/http"
import { getNamespaces } from "./namespaceResource"
import { getDeployments, upDeployment, downDeployment } from "./deploymentResource"
import Refresh from "../../components/Refresh"
import DeploymentActions from "./DeploymentActions"

const INITIAL_STATE = { data: [], loading: false }
const LOADING = "LOADING"
const SET_DATA = "SET_DATA"

function reducer(state, action) {
  switch (action.type) {
    case LOADING:
      return { ...state, loading: true, data: [] }
    case SET_DATA:
      return { ...state, loading: false, data: action.response }
    default:
      return state
  }
}

async function fetchData(dispatch, filter, config) {
  try {
    const result = await getDeployments(filter, config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

async function toUp(context, deployment, setNewPodCount) {
  try {
    setNewPodCount(deployment.name, 1)
    await upDeployment(context, deployment.name, deployment.namespace)
  } catch (e) {
    setNewPodCount(deployment.name, 0)
    notification.error({ message: `Failed to try to up deployment`, description: e.data.error })
  }
}

async function toDown(context, deployment, setNewPodCount) {
  try {
    setNewPodCount(deployment.name, 0)
    await downDeployment(context, deployment.name, deployment.namespace)
  } catch (e) {
    setNewPodCount(deployment.name, 1)
    notification.error({ message: `Failed to try to down deployment`, description: e.data.error })
  }
}

export default function DeploymentPage() {
  const { context } = useParams()
  const [search, setSearch] = useState("")
  const [namespace, setNamespace] = useState("default")
  const [namespaces, setNamespaces] = useState([])
  const [deployments, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    getNamespaces({ context }, { cancelToken: source.token })
      .then((result) => {
        setNamespaces(result.data)
      })
      .catch(() => {})

    return () => {
      source.cancel()
    }
  }, [context])

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { context, namespace }, { cancelToken: source.token })

    return () => {
      source.cancel()
    }
  }, [context, namespace])

  const onReload = useCallback(async () => {
    fetchData(dispatch, { context, namespace })
  }, [context, namespace])

  const updatePodCount = (name, podCount) => {
    const newDeployments = deployments.data.map((dep) =>
      dep.name === name ? { ...dep, pod_count: podCount } : dep,
    )
    dispatch({ type: SET_DATA, response: newDeployments })
  }

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      width: 300,
      sorter: (a, b) => (a.name > b.name) * 2 - 1,
      sortDirections: ["descend", "ascend"],
    },
    {
      title: "Pods Info",
      dataIndex: "pod_info",
      key: "pod_info",
      sorter: (a, b) => (a.pod_info.current > b.pod_info.current) * 2 - 1,
      render: (content) => {
        const color = content.current > 0 ? "green" : "red"
        return (
          <Tag color={color}>
            {content.current}/{content.desired}
          </Tag>
        )
      },
    },
    {
      title: "Action",
      dataIndex: "",
      key: "action",
      width: 140,
      render: (text, deployment) => (
        <DeploymentActions
          context={context}
          deployment={deployment}
          toUp={() => toUp(context, deployment, updatePodCount)}
          toDown={() => toDown(context, deployment, updatePodCount)}
        />
      ),
    },
  ]

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={5} lg={6}>
          <Input.Search
            onChange={(e) => setSearch(e.target.value)}
            onSearch={setSearch}
            value={search}
            enterButton
          />
        </Col>
        <Col xs={6} md={3} xl={2}>
          <Button onClick={() => setSearch("")}>Clear</Button>
        </Col>
        <Col xs={24} md={8} xl={7}>
          <Form.Item label="Namespace">
            <Select
              defaultValue="default"
              value={namespace}
              onChange={setNamespace}
              style={{ width: "100%" }}
            >
              {namespaces.map((ns) => (
                <Select.Option key={ns.name} value={ns.name}>
                  {ns.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Col>
        <Col xs={0} md={8} lg={7} xl={{ span: 6, offset: 3 }} style={{ textAlign: "right" }}>
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          {deployments.data !== [] && (
            <Table
              dataSource={deployments.data.filter(
                (deployment) => search === "" || deployment.name.includes(search),
              )}
              columns={columns}
              rowKey="name"
              loading={deployments.loading}
              size="small"
              scroll={{ x: 600 }}
            />
          )}
        </Col>
      </Row>
    </>
  )
}
