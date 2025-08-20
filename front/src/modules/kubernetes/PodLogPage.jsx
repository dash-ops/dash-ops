import { useCallback, useEffect, useReducer } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router"
import { toast } from "sonner"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { Button } from "@/components/ui/button"
import { ChevronLeft } from "lucide-react"
import { cancelToken } from "../../helpers/http"
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
  } catch {
    toast.error("Ops... Failed to fetch API data")
    dispatch({ type: SET_DATA, response: [] })
  }
}

export default function PodLogPage() {
  const { context } = useParams()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const name = searchParams.get("name") ?? ""
  const namespace = searchParams.get("namespace") ?? "default"
  const [logs, dispatch] = useReducer(reducer, INITIAL_STATE)

  useEffect(() => {
    const source = cancelToken.source()
    dispatch({ type: LOADING })
    fetchData(dispatch, { context, name, namespace }, { cancelToken: source.token })
    return () => {
      source.cancel()
    }
  }, [context, name, namespace])

  const onReload = useCallback(async () => {
    fetchData(dispatch, { context, name, namespace })
  }, [context, name, namespace])

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <Button onClick={() => navigate(-1)} className="gap-2">
          <ChevronLeft className="h-4 w-4" />
          Go Back
        </Button>
        <Refresh onReload={onReload} />
      </div>

      {logs.data.length > 0 && (
        <div className="space-y-4">
          {logs.data.map((l, index) => (
            <Collapsible key={l.name} defaultOpen={index === 0}>
              <CollapsibleTrigger className="flex items-center justify-between w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
                <span className="font-medium">Container: {l.name}</span>
                <ChevronLeft className="h-4 w-4 transition-transform data-[state=open]:rotate-90" />
              </CollapsibleTrigger>
              <CollapsibleContent className="px-4 pb-4">
                <pre className="bg-gray-900 text-green-400 p-4 rounded-md text-sm overflow-auto max-h-96">
                  <code>{l.log}</code>
                </pre>
              </CollapsibleContent>
            </Collapsible>
          ))}
        </div>
      )}

      {logs.loading && (
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900"></div>
          <span className="ml-2">Loading logs...</span>
        </div>
      )}

      {!logs.loading && logs.data.length === 0 && (
        <div className="text-center py-8">
          <span className="text-muted-foreground">No logs found</span>
        </div>
      )}
    </div>
  )
}
