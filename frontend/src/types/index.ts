export interface User {
  id: string
  email: string
  name: string
  organization_id: string
  role: string
  created_at: string
}

export interface Trace {
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

export interface Span {
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
  metadata?: Record<string, any>
}

export interface DashboardData {
  total_traces: number
  total_cost: number
  total_tokens: number
  avg_latency: number
  error_rate: number
  trends: {
    traces: number
    cost: number
    tokens: number
    latency: number
  }
  top_models: Array<{
    model: string
    count: number
    cost: number
  }>
  cost_by_day: Array<{
    date: string
    cost: number
  }>
  traces_by_status: Array<{
    status: string
    count: number
  }>
}

export type TimeRange = '1h' | '24h' | '7d' | '30d' | '90d'
export type TraceStatus = 'success' | 'error' | 'timeout' | 'partial'
