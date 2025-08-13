
import { useState, useEffect, useReducer, useCallback } from "react"
import { useParams, useSearchParams } from "react-router"
import { Row, Col, Table, Tag, notification, Input } from "antd"
import { getNodes } from "./nodesResource"
import Refresh from "../../components/Refresh"
import ProgressData from "./ProgressData"

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

export default function ClusterPage() {
  const { context } = useParams()
  const [searchParams] = useSearchParams()
  const [search, setSearch] = useState(searchParams.get("node") ?? "")
  const [nodes, dispatch] = useReducer(reducer, INITIAL_STATE)

  const fetchData = useCallback(async (config) => {
    try {
      dispatch({ type: LOADING })
      const result = await getNodes({ context }, config)
      dispatch({ type: SET_DATA, response: result.data })
    } catch (e) {
      if (e.message === 'Request canceled') {
        return;
      }
      console.error('Fetch error:', e);
      notification.error({ message: "Ops... Failed to fetch API data" })
      dispatch({ type: SET_DATA, response: [] })
    }
  }, [context])

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;
    fetchData({ signal })

    return () => {
      controller.abort()
    }
  }, [fetchData])

  const onReload = useCallback(() => {
    fetchData()
  }, [fetchData])

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
      render: ({ cpu_requests_fraction }) => <ProgressData percent={cpu_requests_fraction} />,
    },
    {
      title: "CPU limits",
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ cpu_limits_fraction }) => <ProgressData percent={cpu_limits_fraction} />,
    },
    {
      title: "Memory requests",
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ memory_requests_fraction }) => <ProgressData percent={memory_requests_fraction} />,
    },
    {
      title: "Memory limit",
      dataIndex: "allocated_resources",
      key: "allocated_resources",
      render: ({ memory_limits_fraction }) => <ProgressData percent={memory_limits_fraction} />,
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
        <Col xs={24} md={8} xl={7}>
          <Input.Search
            placeholder="Search by node name"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </Col>
        <Col xs={0} md={8} lg={7} xl={{ span: 6, offset: 3 }} style={{ textAlign: "right" }}>
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          <Table
            dataSource={nodes.data.filter((node) =>
              node.name.toLowerCase().includes(search.toLowerCase())
            )}
            columns={columns}
            rowKey="name"
            loading={nodes.loading}
          />
        </Col>
      </Row>
    </>
  )
}
