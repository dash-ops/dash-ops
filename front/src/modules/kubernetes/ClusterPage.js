import React, { useEffect, useReducer } from "react"
import { Row, Col, Table, Tag, Progress, notification } from "antd"
import { cancelToken } from "../../helpers/http"
import { getNodes } from "./nodesResource"
import Refresh from "../../components/Refresh"

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

async function fetchData(dispatch, config) {
  try {
    const result = await getNodes(config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function ClusterPage() {
  const [nodes, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { cancelToken: source.token })

    return () => {
      source.cancel()
    }
  }, [])

  async function onReload() {
    fetchData(dispatch)
  }

  const columns = [
    {
      title: "Node",
      dataIndex: "name",
      key: "name",
      fixed: "left",
      width: 300,
      sorter: (a, b) => (a.name > b.name) * 2 - 1,
      sortDirections: ["descend", "ascend"],
    },
    {
      title: "Ready",
      dataIndex: "ready",
      key: "ready",
      render: (content) => {
        const color = content === "True" ? "green" : "red"
        return <Tag color={color}>{content}</Tag>
      },
    },
    {
      title: "CPU requests",
      dataIndex: "allocatedResources",
      key: "allocatedResources",
      render: ({ cpuRequestsFraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={cpuRequestsFraction.toFixed(1)}
            size="small"
            strokeColor={{
              "0%": "#ffad20",
              "100%": "#ff4d4f",
            }}
            status="active"
          />
        </div>
      ),
    },
    {
      title: "CPU limits",
      dataIndex: "allocatedResources",
      key: "allocatedResources",
      render: ({ cpuLimitsFraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={cpuLimitsFraction.toFixed(1)}
            size="small"
            strokeColor={{
              "0%": "#ffad20",
              "100%": "#ff4d4f",
            }}
            status="active"
          />
        </div>
      ),
    },
    {
      title: "Memory requests",
      dataIndex: "allocatedResources",
      key: "allocatedResources",
      render: ({ memoryRequestsFraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={memoryRequestsFraction.toFixed(1)}
            size="small"
            strokeColor={{
              "0%": "#ffad20",
              "100%": "#ff4d4f",
            }}
            status="active"
          />
        </div>
      ),
    },
    {
      title: "Memory limit",
      dataIndex: "allocatedResources",
      key: "allocatedResources",
      render: ({ memoryLimitsFraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={memoryLimitsFraction.toFixed(1)}
            size="small"
            strokeColor={{
              "0%": "#ffad20",
              "100%": "#ff4d4f",
            }}
            status="active"
          />
        </div>
      ),
    },
    {
      title: "Pods Allocate/Capacity",
      dataIndex: "allocatedResources",
      key: "allocatedResources",
      render: ({ allocatedPods, podCapacity }) => (
        <div style={{ textAlign: "center" }}>
          {allocatedPods}/{podCapacity}
        </div>
      ),
    },
  ]

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={6}></Col>
        <Col xs={6} md={3} xl={2}></Col>
        <Col xs={24} md={6}></Col>
        <Col
          xs={0}
          md={{ span: 6, offset: 3 }}
          xl={{ span: 6, offset: 4 }}
          style={{ textAlign: "right" }}
        >
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          {nodes.data !== [] && (
            <Table
              dataSource={nodes.data}
              columns={columns}
              rowKey="name"
              loading={nodes.loading}
              size="small"
              pagination={false}
              scroll={{ x: 600 }}
            />
          )}
        </Col>
      </Row>
    </>
  )
}
