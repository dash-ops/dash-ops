import React, { useState, useEffect, useReducer } from "react"
import { Link, useLocation, useHistory, useParams } from "react-router-dom"
import { Row, Col, Table, Button, Tooltip, Input, notification, Form, Tag, Select } from "antd"
import { cancelToken } from "../../helpers/http"
import useQuery from "../../helpers/useQuery"
import { getNamespaces } from "./namespaceResource"
import { getPods } from "./podsResource"
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
    const result = await getPods(filter, config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function PodPage() {
  const { context } = useParams()
  const history = useHistory()
  const location = useLocation()
  const query = useQuery()
  const [search, setSearch] = useState(query.get("name") ?? "")
  const [namespace, setNamespace] = useState(query.get("namespace") ?? "default")
  const [namespaces, setNamespaces] = useState([])
  const [pods, dispatch] = useReducer(reducer, INITIAL_STATE)

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

    history.push(`${location.pathname}?name=${search}&namespace=${namespace}`)
    fetchData(dispatch, { context, namespace }, { cancelToken: source.token })

    return () => {
      source.cancel()
    }
  }, [context, namespace]) // eslint-disable-line react-hooks/exhaustive-deps

  async function onReload() {
    fetchData(dispatch, { context, namespace })
  }

  function searchHandler(value) {
    history.push(`${location.pathname}?name=${value}&namespace=${namespace}`)
    setSearch(value)
  }

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      width: 400,
      sorter: (a, b) => (a.name > b.name) * 2 - 1,
      sortDirections: ["descend", "ascend"],
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (content) => {
        const color = content === "True" ? "green" : "red"
        return <Tag color={color}>{content}</Tag>
      },
    },
    {
      title: "Restart",
      dataIndex: "restart_count",
      key: "restart_count",
    },
    {
      title: "Node",
      dataIndex: "node_name",
      key: "node_name",
      width: 300,
      render: (content) => <Link to={`/k8s/${context}?node=${content}`}>{content}</Link>,
    },
    {
      title: "Action",
      dataIndex: "",
      key: "action",
      width: 140,
      render: (text, pod) => (
        <Button.Group>
          <Tooltip title="Pods">
            <Link to={`/k8s/${context}/pod/logs?name=${pod.name}&namespace=${namespace}`}>
              <Button type="primary" ghost size="small">
                Logs
              </Button>
            </Link>
          </Tooltip>
        </Button.Group>
      ),
    },
  ]

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={5} lg={6}>
          <Input.Search
            onChange={(e) => searchHandler(e.target.value)}
            onSearch={searchHandler}
            value={search}
            enterButton
          />
        </Col>
        <Col xs={6} md={3} xl={2}>
          <Button onClick={() => searchHandler("")}>Clear</Button>
        </Col>
        <Col xs={24} md={8} xl={7}>
          <Form.Item label="Namespace">
            <Select
              defaultValue="default"
              value={namespace}
              onChange={setNamespace}
              style={{ width: "100%" }}
            >
              {namespaces.map((ns, index) => (
                <Select.Option key={index} value={ns.name}>
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
          {pods.data !== [] && (
            <Table
              dataSource={pods.data.filter((p) => search === "" || p.name.includes(search))}
              columns={columns}
              rowKey="name"
              loading={pods.loading}
              size="small"
              scroll={{ x: 600 }}
            />
          )}
        </Col>
      </Row>
    </>
  )
}
