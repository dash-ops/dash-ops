
import { useState, useEffect, useReducer, useCallback } from "react"
import { useParams, useSearchParams } from "react-router"
import { toast } from "sonner"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
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
      toast.error("Ops... Failed to fetch API data")
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

  const filteredData = nodes.data.filter((node) =>
    node.name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="md:col-span-1">
          <Input
            placeholder="Search by node name"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="hidden md:block" />
        <div className="flex justify-end">
          <Refresh onReload={onReload} />
        </div>
      </div>
      
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[300px]">Node</TableHead>
              <TableHead>Ready</TableHead>
              <TableHead>CPU requests</TableHead>
              <TableHead>CPU limits</TableHead>
              <TableHead>Memory requests</TableHead>
              <TableHead>Memory limit</TableHead>
              <TableHead className="text-center">Pods Allocate/Capacity</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {nodes.loading ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-8">
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900"></div>
                    <span className="ml-2">Loading...</span>
                  </div>
                </TableCell>
              </TableRow>
            ) : filteredData.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-8">
                  <span className="text-muted-foreground">No nodes found</span>
                </TableCell>
              </TableRow>
            ) : (
              filteredData.map((node) => (
                <TableRow key={node.name}>
                  <TableCell className="font-medium">{node.name}</TableCell>
                  <TableCell>
                    <Badge variant={node.ready === "True" ? "default" : "destructive"}>
                      {node.ready}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <ProgressData percent={node.allocated_resources?.cpu_requests_fraction} />
                  </TableCell>
                  <TableCell>
                    <ProgressData percent={node.allocated_resources?.cpu_limits_fraction} />
                  </TableCell>
                  <TableCell>
                    <ProgressData percent={node.allocated_resources?.memory_requests_fraction} />
                  </TableCell>
                  <TableCell>
                    <ProgressData percent={node.allocated_resources?.memory_limits_fraction} />
                  </TableCell>
                  <TableCell className="text-center">
                    {node.allocated_resources?.allocated_pods}/{node.allocated_resources?.pod_capacity}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
