import React, { useEffect, useReducer } from "react"
import { useParams, useHistory } from "react-router-dom"
import { Row, Col, Collapse, Button, notification } from "antd"
import { CaretLeftOutlined } from "@ant-design/icons"
import { cancelToken } from "../../helpers/http"
import useQuery from "../../helpers/useQuery"
import { getPodLogs } from "./podsResource"
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
    const result = await getPodLogs(filter, config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (e) {
    notification.error({ message: "Ops... Failed to fetch API data" })
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function PodLogPage() {
  const { context } = useParams()
  const history = useHistory()
  const query = useQuery()
  const name = query.get("name") ?? ""
  const namespace = query.get("namespace") ?? "default"
  const [logs, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { context, name, namespace }, { cancelToken: source.token })
    return () => {
      source.cancel()
    }
  }, [context, name, namespace])

  async function onReload() {
    fetchData(dispatch, { context, name, namespace })
  }

  return (
    <>
      <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
        <Col xs={18} md={5} lg={6}>
          <Button type="primary" icon={<CaretLeftOutlined />} onClick={history.goBack}>
            Go Back
          </Button>
        </Col>
        <Col xs={6} md={3} xl={2}/>
        <Col xs={24} md={8} xl={7}/>
        <Col xs={0} md={8} lg={7} xl={{ span: 6, offset: 3 }} style={{ textAlign: "right" }}>
          <Refresh onReload={onReload} />
        </Col>
      </Row>
      <Row>
        <Col flex="auto" style={{ marginTop: 10 }}>
          {logs.data !== [] && (
            <Collapse defaultActiveKey={["0"]}>
              {logs.data.map((l, index) => (
                <Collapse.Panel header={`Container: ${l.name}`} key={index}>
                  <pre>
                    <code>{l.log}</code>
                  </pre>
                </Collapse.Panel>
              ))}
            </Collapse>
          )}
        </Col>
      </Row>
    </>
  )
}
