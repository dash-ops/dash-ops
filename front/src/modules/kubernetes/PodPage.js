import React, { useState, useEffect, useReducer } from "react"
import { Row, Col, Table, Button, Input, notification, Form, Tag, Select } from "antd"
import { cancelToken } from "../../helpers/http"
import { getNamespaces } from "./namespaceResource"
import { getPods } from "./podsResource"
import Refresh from "../../components/Refresh"
import { Link } from "react-router-dom"

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
  const [search, setSearch] = useState("")
  const [namespace, setNamespace] = useState("default")
  const [namespaces, setNamespaces] = useState([])
  const [pods, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    getNamespaces({ cancelToken: source.token })
      .then((result) => {
        setNamespaces(result.data)
      })
      .catch(() => {})

    return () => {
      source.cancel()
    }
  }, [])

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { namespace }, { cancelToken: source.token })

    return () => {
      source.cancel()
    }
  }, [namespace])

  async function onReload() {
    fetchData(dispatch)
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
      render: (content) => <Link to={`/k8s?node=${content}`}>{content}</Link>,
    },
  ]

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={6}>
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
        <Col xs={24} md={6}>
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
          {pods.data !== [] && (
            <Table
              dataSource={pods.data}
              columns={columns}
              rowKey="name"
              loading={pods.loading}
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
