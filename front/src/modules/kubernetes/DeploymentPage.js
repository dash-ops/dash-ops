import React, { useState, useEffect, useReducer } from "react"
import { Row, Col, Table, Button, Input, notification, Tag } from "antd"
import { getDeployments, upDeployment, downDeployment } from "./deploymentResource"
import Refresh from "../../components/Refresh"
import DeploymentActions from "./DeploymentActions"

const initialState = { data: [], loading: false }
const LOADING = "LOADING"
const SET_DATA = "SET_DATA"

function reducer(state, action) {
  switch (action.type) {
    case LOADING:
      return { loading: true, data: [] }
    case SET_DATA:
      return { loading: false, data: action.response }
    default:
      return initialState
  }
}

function useSearchInput(initalValue) {
  const [value, setValue] = useState(initalValue)

  return {
    value,
    onSearch: v => setValue(v),
  }
}

async function fetchData(dispatch) {
  try {
    const result = await getDeployments()
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

async function toUp(deployment, setNewPodCount) {
  try {
    await upDeployment(deployment.name, deployment.namespace)
    setNewPodCount(deployment.name, 1)
  } catch (e) {
    notification.error({ message: `Failed to try to up deployment` })
  }
}

async function toDown(deployment, setNewPodCount) {
  try {
    await downDeployment(deployment.name, deployment.namespace)
    setNewPodCount(deployment.name, 0)
  } catch (e) {
    notification.error({ message: `Failed to try to down deployment` })
  }
}

export default function DeploymentPage() {
  const search = useSearchInput("")
  const [deployments, dispatch] = useReducer(reducer, initialState)
  useEffect(() => {
    dispatch({ type: LOADING })
    fetchData(dispatch)
  }, [])

  async function onReload() {
    fetchData(dispatch)
  }

  function updatePodCount(name, podCount) {
    const newDeployments = deployments.data.map(dep =>
      dep.name === name ? { ...dep, pod_count: podCount } : dep,
    )
    dispatch({ type: SET_DATA, response: newDeployments })
  }

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      fixed: "left",
      width: 300,
      sorter: (a, b) => (a.name > b.name) * 2 - 1,
      sortDirections: ["descend", "ascend"],
    },
    { title: "Namespace", dataIndex: "namespace", key: "namespace" },
    {
      title: "Pods running",
      dataIndex: "pod_count",
      key: "pod_count",
      render: content => {
        const color = content > 0 ? "green" : "red"
        return <Tag color={color}>{content}</Tag>
      },
    },
    {
      title: "Action",
      dataIndex: "",
      key: "action",
      fixed: "right",
      width: 90,
      render: (text, deployment) => (
        <DeploymentActions
          deployment={deployment}
          toUp={() => toUp(deployment, updatePodCount)}
          toDown={() => toDown(deployment, updatePodCount)}
        />
      ),
    },
  ]

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={6}>
          <Input.Search
            onChange={e => search.onSearch(e.target.value)}
            onSearch={search.onSearch}
            value={search.value}
            enterButton
          />
        </Col>
        <Col xs={6} md={6}>
          <Button onClick={() => search.onSearch("")}>Clear</Button>
        </Col>
        <Col xs={0} md={{ span: 6, offset: 6 }} style={{ textAlign: "right" }}>
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          <Table
            dataSource={deployments.data.filter(
              deployment => search.value === "" || deployment.name.includes(search.value),
            )}
            columns={columns}
            rowKey="name"
            loading={deployments.loading}
            size="small"
            pagination={false}
            scroll={{ x: 600 }}
          />
        </Col>
      </Row>
    </>
  )
}