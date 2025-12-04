import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { Badge } from '../components/ui/Badge'
import { 
  ArrowLeft,
  Clock,
  DollarSign,
  Zap,
  Activity,
  Code,
  AlertCircle,
  CheckCircle2
} from 'lucide-react'
import api from '../services/api'
import { formatCurrency, formatDuration, formatDate, getStatusColor } from '../lib/utils'

interface Span {
  span_id: string
  trace_id: string
  parent_span_id?: string
  name: string
  start_time: string
  end_time: string
  duration_ms: number
  model: string
  provider: string
  input?: string
  output?: string
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
  cost_usd: number
  status: string
  error_message?: string
}

interface TraceDetail {
  trace_id: string
  organization_id: string
  project_id: string
  timestamp: string
  trace_type: string
  duration_ms: number
  status: string
  total_cost_usd: number
  total_tokens: number
  model: string
  provider: string
  user_id?: string
  metadata?: Record<string, any>
  spans: Span[]
}

export default function TraceDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [trace, setTrace] = useState<TraceDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [selectedSpan, setSelectedSpan] = useState<Span | null>(null)

  useEffect(() => {
    if (id) {
      loadTrace(id)
    }
  }, [id])

  const loadTrace = async (traceId: string) => {
    setLoading(true)
    try {
      const response = await api.getTrace(traceId)
      console.log('Trace detail:', response)
      setTrace(response.data)
      if (response.data.spans && response.data.spans.length > 0) {
        setSelectedSpan(response.data.spans[0])
      }
    } catch (error) {
      console.error('Failed to load trace:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading trace...</p>
        </div>
      </div>
    )
  }

  if (!trace) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <p className="text-lg font-medium">Trace not found</p>
          <button
            onClick={() => navigate('/traces')}
            className="mt-4 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
          >
            Back to Traces
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate('/traces')}
          className="p-2 hover:bg-accent rounded-md"
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-bold">Trace Details</h1>
            <Badge 
              variant={trace.status === 'success' ? 'success' : 'destructive'}
              className={getStatusColor(trace.status)}
            >
              {trace.status}
            </Badge>
          </div>
          <p className="text-muted-foreground mt-2 font-mono text-sm">
            {trace.trace_id}
          </p>
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Duration</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatDuration(trace.duration_ms)}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cost</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(trace.total_cost_usd)}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Tokens</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{trace.total_tokens.toLocaleString()}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Spans</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{trace.spans?.length || 0}</div>
          </CardContent>
        </Card>
      </div>

      {/* Trace Info */}
      <Card>
        <CardHeader>
          <CardTitle>Trace Information</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="text-sm font-medium text-muted-foreground">Model</label>
              <p className="text-lg font-medium mt-1">{trace.model}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Provider</label>
              <p className="text-lg font-medium mt-1 capitalize">{trace.provider}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Type</label>
              <p className="text-lg font-medium mt-1">{trace.trace_type}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Timestamp</label>
              <p className="text-lg font-medium mt-1">{formatDate(trace.timestamp)}</p>
            </div>
            {trace.user_id && (
              <div>
                <label className="text-sm font-medium text-muted-foreground">User ID</label>
                <p className="text-lg font-medium mt-1 font-mono text-sm">{trace.user_id}</p>
              </div>
            )}
            {trace.metadata && Object.keys(trace.metadata).length > 0 && (
              <div className="md:col-span-2">
                <label className="text-sm font-medium text-muted-foreground">Metadata</label>
                <pre className="text-sm mt-1 p-3 bg-muted rounded-md overflow-x-auto">
                  {JSON.stringify(trace.metadata, null, 2)}
                </pre>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Spans */}
      {trace.spans && trace.spans.length > 0 && (
        <div className="grid gap-4 md:grid-cols-3">
          {/* Spans List */}
          <Card className="md:col-span-1">
            <CardHeader>
              <CardTitle>Spans ({trace.spans.length})</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {trace.spans.map((span, index) => (
                  <div
                    key={span.span_id}
                    onClick={() => setSelectedSpan(span)}
                    className={`p-3 rounded-lg border cursor-pointer transition-colors ${
                      selectedSpan?.span_id === span.span_id
                        ? 'border-primary bg-primary/5'
                        : 'hover:bg-muted/50'
                    }`}
                  >
                    <div className="flex items-start justify-between gap-2">
                      <div className="flex-1 min-w-0">
                        <div className="font-medium text-sm truncate">
                          {index + 1}. {span.name}
                        </div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {formatDuration(span.duration_ms)} • {span.model}
                        </div>
                      </div>
                      {span.status === 'success' ? (
                        <CheckCircle2 className="h-4 w-4 text-green-600 flex-shrink-0" />
                      ) : (
                        <AlertCircle className="h-4 w-4 text-red-600 flex-shrink-0" />
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Span Details */}
          <Card className="md:col-span-2">
            <CardHeader>
              <CardTitle>Span Details</CardTitle>
            </CardHeader>
            <CardContent>
              {selectedSpan ? (
                <div className="space-y-4">
                  {/* Span Stats */}
                  <div className="grid gap-4 md:grid-cols-3">
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">Duration</label>
                      <p className="text-lg font-medium mt-1">{formatDuration(selectedSpan.duration_ms)}</p>
                    </div>
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">Cost</label>
                      <p className="text-lg font-medium mt-1">{formatCurrency(selectedSpan.cost_usd)}</p>
                    </div>
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">Tokens</label>
                      <p className="text-lg font-medium mt-1">
                        {selectedSpan.prompt_tokens} + {selectedSpan.completion_tokens} = {selectedSpan.total_tokens}
                      </p>
                    </div>
                  </div>

                  {/* Input/Output */}
                  {selectedSpan.input && (
                    <div>
                      <label className="text-sm font-medium text-muted-foreground flex items-center gap-2 mb-2">
                        <Code className="h-4 w-4" />
                        Input
                      </label>
                      <div className="p-3 bg-muted rounded-md text-sm max-h-48 overflow-y-auto">
                        {selectedSpan.input}
                      </div>
                    </div>
                  )}

                  {selectedSpan.output && (
                    <div>
                      <label className="text-sm font-medium text-muted-foreground flex items-center gap-2 mb-2">
                        <Code className="h-4 w-4" />
                        Output
                      </label>
                      <div className="p-3 bg-muted rounded-md text-sm max-h-48 overflow-y-auto">
                        {selectedSpan.output}
                      </div>
                    </div>
                  )}

                  {selectedSpan.error_message && (
                    <div>
                      <label className="text-sm font-medium text-red-600 flex items-center gap-2 mb-2">
                        <AlertCircle className="h-4 w-4" />
                        Error Message
                      </label>
                      <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md text-sm">
                        {selectedSpan.error_message}
                      </div>
                    </div>
                  )}

                  {/* Span Metadata */}
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">Span Details</label>
                    <div className="mt-2 grid gap-2 text-sm">
                      <div className="flex justify-between p-2 rounded bg-muted/50">
                        <span className="text-muted-foreground">Span ID:</span>
                        <span className="font-mono">{selectedSpan.span_id.substring(0, 16)}...</span>
                      </div>
                      <div className="flex justify-between p-2 rounded bg-muted/50">
                        <span className="text-muted-foreground">Model:</span>
                        <span className="font-medium">{selectedSpan.model}</span>
                      </div>
                      <div className="flex justify-between p-2 rounded bg-muted/50">
                        <span className="text-muted-foreground">Provider:</span>
                        <span className="font-medium capitalize">{selectedSpan.provider}</span>
                      </div>
                      <div className="flex justify-between p-2 rounded bg-muted/50">
                        <span className="text-muted-foreground">Start Time:</span>
                        <span>{formatDate(selectedSpan.start_time)}</span>
                      </div>
                      <div className="flex justify-between p-2 rounded bg-muted/50">
                        <span className="text-muted-foreground">End Time:</span>
                        <span>{formatDate(selectedSpan.end_time)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-12 text-muted-foreground">
                  Select a span to view details
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Span Timeline */}
      {trace.spans && trace.spans.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Span Timeline</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {trace.spans.map((span, index) => {
                const startTime = new Date(span.start_time).getTime()
                const traceStartTime = new Date(trace.timestamp).getTime()
                const offset = startTime - traceStartTime
                const totalDuration = trace.duration_ms
                const startPercent = (offset / totalDuration) * 100
                const widthPercent = (span.duration_ms / totalDuration) * 100

                return (
                  <div key={span.span_id} className="space-y-1">
                    <div className="flex items-center justify-between text-sm">
                      <span className="font-medium">{index + 1}. {span.name}</span>
                      <span className="text-muted-foreground">{formatDuration(span.duration_ms)}</span>
                    </div>
                    <div className="relative h-8 bg-muted rounded-md overflow-hidden">
                      <div
                        className={`absolute h-full rounded transition-all ${
                          span.status === 'success' 
                            ? 'bg-green-500' 
                            : span.status === 'error'
                            ? 'bg-red-500'
                            : 'bg-yellow-500'
                        }`}
                        style={{
                          left: `${Math.max(0, startPercent)}%`,
                          width: `${widthPercent}%`,
                        }}
                      >
                        <div className="flex items-center justify-center h-full text-xs text-white font-medium">
                          {span.model}
                        </div>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}