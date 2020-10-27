import { useState, useEffect, useReducer } from "react";
import { useParams } from "react-router-dom"
import { Row, Col, Table, Tag, Progress, notification, Input, Button } from "antd"
import { cancelToken } from "../../helpers/http"
import useQuery from "../../helpers/useQuery"
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

async function fetchData(dispatch, filter, config) {
  try {
    const result = await getNodes(filter, config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function ClusterPage() {
  const { context } = useParams()
  const query = useQuery()
  const [search, setSearch] = useState(query.get("node") ?? "")
  const [nodes, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { context }, { cancelToken: source.token })

    return () => {
      source.cancel()
    }
  }, [context])

  async function onReload() {
    fetchData(dispatch, { context })
  }

  const columns = [
    {
      title: "Node",
      dataIndex: "name",
      key: "name",
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
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ cpu_requests_fraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={cpu_requests_fraction.toFixed(1)}
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
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ cpu_limits_fraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={cpu_limits_fraction.toFixed(1)}
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
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ memory_requests_fraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={memory_requests_fraction.toFixed(1)}
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
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ memory_limits_fraction }) => (
        <div style={{ maxWidth: 170 }}>
          <Progress
            percent={memory_limits_fraction.toFixed(1)}
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
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ allocated_pods, pod_capacity }) => (
        <div style={{ textAlign: "center" }}>
          {allocated_pods}/{pod_capacity}
        </div>
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
        <Col xs={24} md={8} xl={7}></Col>
        <Col xs={0} md={8} lg={7} xl={{ span: 6, offset: 3 }} style={{ textAlign: "right" }}>
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          {nodes.data !== [] && (
            <Table
              dataSource={nodes.data.filter((n) => search === "" || n.name.includes(search))}
              columns={columns}
              rowKey="name"
              loading={nodes.loading}
              size="small"
              scroll={{ x: 600 }}
            />
          )}
        </Col>
      </Row>
    </>
  )
}
