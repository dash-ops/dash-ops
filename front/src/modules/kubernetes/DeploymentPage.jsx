import { useState, useEffect, useReducer, useCallback } from "react"
import { useParams } from "react-router"
import { toast } from "sonner"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { getNamespaces } from "./namespaceResource"
import { getDeployments, upDeployment, downDeployment } from "./deploymentResource"
import Refresh from "../../components/Refresh"
import DeploymentActions from "./DeploymentActions"

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
    const result = await getDeployments(filter, config)
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

async function toUp(context, deployment, setNewPodCount) {
  try {
    setNewPodCount(deployment.name, 1)
    await upDeployment(context, deployment.name, deployment.namespace)
  } catch (e) {
    setNewPodCount(deployment.name, 0)
    toast.error(`Failed to try to up deployment: ${e.data.error}`)
  }
}

async function toDown(context, deployment, setNewPodCount) {
  try {
    setNewPodCount(deployment.name, 0)
    await downDeployment(context, deployment.name, deployment.namespace)
  } catch (e) {
    setNewPodCount(deployment.name, 1)
    toast.error(`Failed to try to down deployment: ${e.data.error}`)
  }
}

export default function DeploymentPage() {
  const { context } = useParams()
  const [search, setSearch] = useState("")
  const [namespace, setNamespace] = useState("default")
  const [namespaces, setNamespaces] = useState([])
  const [deployments, dispatch] = useReducer(reducer, INITIAL_STATE)

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

  const handleNamespaceChange = (newNamespace) => {
    setNamespace(newNamespace)
  }

  const updatePodCount = (name, podCount) => {
    const newDeployments = deployments.data.map((dep) =>
      dep.name === name ? { ...dep, pod_count: podCount } : dep,
    )
    dispatch({ type: SET_DATA, response: newDeployments })
  }

  const filteredData = deployments.data.filter(
    (deployment) => search === "" || deployment.name.includes(search)
  )

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-12 gap-4">
        <div className="md:col-span-3">
          <Input
            placeholder="Search deployments..."
            onChange={(e) => setSearch(e.target.value)}
            value={search}
          />
        </div>
        <div className="md:col-span-1">
          <Button 
            variant="outline" 
            onClick={() => setSearch("")}
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

      {deployments.data.length > 0 && (
        <div className="border rounded-lg overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[300px]">Name</TableHead>
                <TableHead>Pods Info</TableHead>
                <TableHead className="w-[140px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {deployments.loading ? (
                <TableRow>
                  <TableCell colSpan={3} className="text-center py-8">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900"></div>
                      <span className="ml-2">Loading...</span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : filteredData.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={3} className="text-center py-8">
                    <span className="text-muted-foreground">No deployments found</span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((deployment) => (
                  <TableRow key={deployment.name}>
                    <TableCell className="font-medium">{deployment.name}</TableCell>
                    <TableCell>
                      <Badge 
                        variant={deployment.pod_info.current > 0 ? "default" : "destructive"}
                      >
                        {deployment.pod_info.current}/{deployment.pod_info.desired}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <DeploymentActions
                        context={context}
                        deployment={deployment}
                        toUp={() => toUp(context, deployment, updatePodCount)}
                        toDown={() => toDown(context, deployment, updatePodCount)}
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
