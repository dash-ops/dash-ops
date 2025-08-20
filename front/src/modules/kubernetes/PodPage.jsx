import { useState, useEffect, useReducer, useCallback } from "react"
import { Link, useLocation, useNavigate, useParams, useSearchParams } from "react-router"
import { toast } from "sonner"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import { cancelToken } from "../../helpers/http" // Deprecated - using AbortController instead
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
    if (e.message === 'Request canceled') {
      return;
    }
    console.error('Fetch error:', e);
    toast.error("Ops... Failed to fetch API data")
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function PodPage() {
  const { context } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const [search, setSearch] = useState(searchParams.get("name") ?? "")
  const [namespace, setNamespace] = useState(searchParams.get("namespace") ?? "default")
  const [namespaces, setNamespaces] = useState([])
  const [pods, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const controller = new AbortController()
    const signal = controller.signal
    
    getNamespaces({ context }, { signal })
      .then((result) => {
        setNamespaces(result.data)
      })
      .catch((e) => {
        if (e.message !== 'Request canceled') {
          console.error('Error fetching namespaces:', e);
        }
      })

    return () => {
      controller.abort()
    }
  }, [context])

  useEffect(() => {
    const controller = new AbortController()
    const signal = controller.signal
    dispatch({ type: LOADING })

    fetchData(dispatch, { context, namespace }, { signal })

    return () => {
      controller.abort()
    }
  }, [context, namespace])

  const onReload = useCallback(async () => {
    fetchData(dispatch, { context, namespace })
  }, [context, namespace])

  const searchHandler = (value) => {
    setSearch(value)
    navigate(`${location.pathname}?name=${value}&namespace=${namespace}`)
  }

  const handleNamespaceChange = (newNamespace) => {
    setNamespace(newNamespace)
    navigate(`${location.pathname}?name=${search}&namespace=${newNamespace}`)
  }

  const filteredData = pods.data.filter((p) => search === "" || p.name.includes(search))

  const getStatusColor = (status) => {
    switch (status) {
      case "Running":
        return "default"
      case "Succeeded":
        return "secondary"
      case "Pending":
        return "outline"
      default:
        return "destructive"
    }
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-12 gap-4">
        <div className="md:col-span-3">
          <Input
            placeholder="Search pods..."
            onChange={(e) => searchHandler(e.target.value)}
            value={search}
          />
        </div>
        <div className="md:col-span-1">
          <Button 
            variant="outline" 
            onClick={() => searchHandler("")}
            className="w-full"
          >
            Clear
          </Button>
        </div>
        <div className="md:col-span-3">
          <div className="space-y-1">
            <Select value={namespace} onValueChange={handleNamespaceChange}>
              <SelectTrigger>
                <SelectValue placeholder="Select namespace" />
              </SelectTrigger>
              <SelectContent>
                {namespaces.map((ns) => (
                  <SelectItem key={ns.name} value={ns.name}>
                    {ns.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <div className="hidden md:block md:col-span-2" />
        <div className="md:col-span-3 flex justify-end">
          <Refresh onReload={onReload} />
        </div>
      </div>

      {pods.data.length > 0 && (
        <div className="border rounded-lg overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[400px]">Name</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Restart</TableHead>
                <TableHead className="w-[140px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pods.loading ? (
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
                    <span className="text-muted-foreground">No pods found</span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((pod) => (
                  <TableRow key={pod.name}>
                    <TableCell className="font-medium">{pod.name}</TableCell>
                    <TableCell>
                      <Badge variant={getStatusColor(pod.condition_status.status)}>
                        {pod.condition_status.status}
                      </Badge>
                    </TableCell>
                    <TableCell>{pod.restart_count}</TableCell>
                    <TableCell className="text-right">
                      <TooltipProvider>
                        <div className="flex gap-1 justify-end">
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="outline" size="sm" asChild>
                                <Link to={`/k8s/${context}/pod/logs?name=${pod.name}&namespace=${namespace}`}>
                                  Logs
                                </Link>
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              <p>Containers Log</p>
                            </TooltipContent>
                          </Tooltip>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="outline" size="sm" asChild>
                                <Link to={`/k8s/${context}?node=${pod.node_name}`}>
                                  Node
                                </Link>
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              <p>Details {pod.node_name}</p>
                            </TooltipContent>
                          </Tooltip>
                        </div>
                      </TooltipProvider>
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
