import { useState, useEffect, useReducer, useCallback } from "react"
import { useParams } from "react-router"
import { toast } from "sonner"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
// import { cancelToken } from "../../helpers/http" // Deprecated - using AbortController instead
import { getInstances, startInstance, stopInstance } from "./instanceResource"
import Refresh from "../../components/Refresh"
import InstanceActions from "./InstanceActions"
import InstanceTag from "./InstanceTag"

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
    const result = await getInstances(filter, config)
    dispatch({ type: SET_DATA, response: result.data })
  } catch (error) {
    if (error.message === 'Request canceled') {
      return;
    }
    console.error('Fetch error:', error);
    toast.error("Ops... Failed to fetch API data")
    dispatch({ type: SET_DATA, response: [] })
  }
}

async function toStart(key, instance, setNewState) {
  try {
    setNewState(instance.instance_id, "pending")
    const response = await startInstance(key, instance.instance_id)
    setNewState(instance.instance_id, response.data.current_state)
  } catch {
    setNewState(instance.instance_id, "stopped")
    toast.error("Failed to try to start Instance")
  }
}

async function toStop(key, instance, setNewState) {
  try {
    setNewState(instance.instance_id, "stopping")
    const response = await stopInstance(key, instance.instance_id)
    setNewState(instance.instance_id, response.data.current_state)
  } catch {
    setNewState(instance.instance_id, "running")
    toast.error("Failed to try to stop Instance")
  }
}

export default function InstancePage() {
  const { key } = useParams()
  const [search, setSearch] = useState("")
  const [instances, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const controller = new AbortController()
    const signal = controller.signal
    dispatch({ type: LOADING })
    fetchData(dispatch, { accountKey: key }, { signal })

    return () => {
      controller.abort()
    }
  }, [key])

  const onReload = useCallback(async () => {
    fetchData(dispatch, { accountKey: key })
  }, [key])

  const updateInstanceState = (id, state) => {
    const newInstances = instances.data.map((inst) =>
      inst.instance_id === id ? { ...inst, state } : inst,
    )
    dispatch({ type: SET_DATA, response: newInstances })
  }

  const filteredData = instances.data.filter(
    (instance) => search === "" || instance.name.includes(search)
  )

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="md:col-span-1">
          <Input
            placeholder="Search instances..."
            onChange={(e) => setSearch(e.target.value)}
            value={search}
          />
        </div>
        <div>
          <Button 
            variant="outline" 
            onClick={() => setSearch("")}
            className="w-full md:w-auto"
          >
            Clear
          </Button>
        </div>
        <div className="hidden md:block" />
        <div className="flex justify-end">
          <Refresh onReload={onReload} />
        </div>
      </div>

      {instances.data.length > 0 && (
        <div className="border rounded-lg overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[300px]">Instance</TableHead>
                <TableHead>Instance Id</TableHead>
                <TableHead>State</TableHead>
                <TableHead className="w-[120px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {instances.loading ? (
                <TableRow>
                  <TableCell colSpan={4} className="text-center py-8">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900"></div>
                      <span className="ml-2">Loading...</span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : filteredData.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} className="text-center py-8">
                    <span className="text-muted-foreground">No instances found</span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((instance) => (
                  <TableRow key={instance.instance_id}>
                    <TableCell className="font-medium">{instance.name}</TableCell>
                    <TableCell className="text-sm font-mono">{instance.instance_id}</TableCell>
                    <TableCell>
                      {instance.state && <InstanceTag state={instance.state} />}
                    </TableCell>
                    <TableCell className="text-right">
                      <InstanceActions
                        instance={instance}
                        toStart={() => toStart(key, instance, updateInstanceState)}
                        toStop={() => toStop(key, instance, updateInstanceState)}
                      />
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
