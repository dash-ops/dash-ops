import React, { useState, useEffect, useReducer } from "react"
import { Row, Col, Table, Button, Input, notification } from "antd"
import { getInstances, startInstance, stopInstance } from "./instanceResource"
import Refresh from "../../components/Refresh"
import InstanceActions from "./InstanceActions"
import InstanceTag from "./InstanceTag"

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
    const result = await getInstances()
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

async function toStart(instance, setNewState) {
  try {
    const response = await startInstance(instance.instance_id)
    setNewState(instance.instance_id, response.data.current_state)
  } catch (e) {
    notification.error({ message: "Failed to try to start Instance" })
  }
}

async function toStop(instance, setNewState) {
  try {
    const response = await stopInstance(instance.instance_id)
    setNewState(instance.instance_id, response.data.current_state)
  } catch (e) {
    notification.error({ message: "Failed to try to stop Instance" })
  }
}

export default function InstancePage() {
  const search = useSearchInput("")
  const [instances, dispatch] = useReducer(reducer, initialState)
  useEffect(() => {
    dispatch({ type: LOADING })
    fetchData(dispatch)
  }, [])

  async function onReload() {
    fetchData(dispatch)
  }

  function updateInstanceState(id, state) {
    const newInstances = instances.data.map(inst =>
      inst.instance_id === id ? { ...inst, state } : inst,
    )
    dispatch({ type: SET_DATA, response: newInstances })
  }

  const columns = [
    {
      title: "Instance",
      dataIndex: "name",
      key: "name",
      fixed: "left",
      width: 300,
      sorter: (a, b) => (a.name > b.name) * 2 - 1,
      sortDirections: ["descend", "ascend"],
    },
    { title: "Instance Id", dataIndex: "instance_id", key: "instance_id" },
    {
      title: "State",
      dataIndex: "state",
      key: "state",
      render: state => !state || <InstanceTag state={state} />,
    },
    {
      title: "Action",
      dataIndex: "",
      key: "action",
      fixed: "right",
      width: 120,
      render: (text, instance) => (
        <InstanceActions
          instance={instance}
          toStart={() => toStart(instance, updateInstanceState)}
          toStop={() => toStop(instance, updateInstanceState)}
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
            dataSource={instances.data.filter(
              instance => search.value === "" || instance.name.includes(search.value),
            )}
            columns={columns}
            rowKey="instance_id"
            loading={instances.loading}
            size="small"
            pagination={false}
            scroll={{ x: 600 }}
          />
        </Col>
      </Row>
    </>
  )
}
