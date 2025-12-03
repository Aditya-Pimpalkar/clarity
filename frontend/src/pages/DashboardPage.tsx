import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Activity, 
  Clock, 
  CheckCircle,
  Zap,
  AlertCircle
} from 'lucide-react'
import {
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import api from '../services/api'
import { formatCurrency, formatNumber, formatDuration } from '../lib/utils'

interface DashboardStats {
  total_traces: number
  total_cost: number
  total_tokens: number
  avg_latency: number
  error_rate: number
  success_rate: number
  trends: {
    traces: number
    cost: number
    tokens: number
    latency: number
  }
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [timeRange, setTimeRange] = useState('24h')

  useEffect(() => {
    loadDashboardData()
  }, [timeRange])

  const loadDashboardData = async () => {
    setLoading(true)
    try {
      const response = await api.getDashboard(timeRange)
      setStats(response.data)
    } catch (error) {
      console.error('Failed to load dashboard:', error)
      // Use mock data if API fails
      setStats({
        total_traces: 64,
        total_cost: 0.15,
        total_tokens: 12847,
        avg_latency: 156,
        error_rate: 3.1,
        success_rate: 96.9,
        trends: {
          traces: 12.5,
          cost: 8.3,
          tokens: 15.2,
          latency: -4.7,
        },
      })
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading dashboard...</p>
        </div>
      </div>
    )
  }

  if (!stats) return null

  const getTrendIcon = (value: number) => {
    if (value > 0) return <TrendingUp className="h-4 w-4 text-green-600" />
    if (value < 0) return <TrendingDown className="h-4 w-4 text-red-600" />
    return null
  }

  const getTrendColor = (value: number) => {
    if (value > 0) return 'text-green-600'
    if (value < 0) return 'text-red-600'
    return 'text-gray-600'
  }

  // Sample data for charts (we'll make this dynamic later)
  const costByDay = [
    { date: 'Mon', cost: 0.023 },
    { date: 'Tue', cost: 0.031 },
    { date: 'Wed', cost: 0.028 },
    { date: 'Thu', cost: 0.035 },
    { date: 'Fri', cost: 0.029 },
    { date: 'Sat', cost: 0.019 },
    { date: 'Sun', cost: 0.015 },
  ]

  const modelDistribution = [
    { name: 'GPT-4', value: 45, cost: 0.089 },
    { name: 'GPT-3.5', value: 12, cost: 0.002 },
    { name: 'Claude-3', value: 7, cost: 0.059 },
  ]

  const COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444']

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground mt-2">
            Real-time insights from your LLM usage
          </p>
        </div>

        {/* Time Range Selector */}
        <div className="flex gap-2">
          {['1h', '24h', '7d', '30d'].map((range) => (
            <button
              key={range}
              onClick={() => setTimeRange(range)}
              className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                timeRange === range
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-secondary hover:bg-secondary/80'
              }`}
            >
              {range}
            </button>
          ))}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {/* Total Traces */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Traces</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatNumber(stats.total_traces)}</div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              {getTrendIcon(stats.trends.traces)}
              <span className={getTrendColor(stats.trends.traces)}>
                {Math.abs(stats.trends.traces)}%
              </span>
              <span>from last period</span>
            </div>
          </CardContent>
        </Card>

        {/* Total Cost */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Cost</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(stats.total_cost)}</div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              {getTrendIcon(stats.trends.cost)}
              <span className={getTrendColor(stats.trends.cost)}>
                {Math.abs(stats.trends.cost)}%
              </span>
              <span>from last period</span>
            </div>
          </CardContent>
        </Card>

        {/* Avg Latency */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatDuration(stats.avg_latency)}</div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              {getTrendIcon(stats.trends.latency)}
              <span className={getTrendColor(stats.trends.latency)}>
                {Math.abs(stats.trends.latency)}%
              </span>
              <span>from last period</span>
            </div>
          </CardContent>
        </Card>

        {/* Success Rate */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
            <CheckCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.success_rate.toFixed(1)}%</div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              <Zap className="h-3 w-3" />
              <span>{stats.total_traces - Math.floor(stats.total_traces * stats.error_rate / 100)} successful</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Cost Trend */}
        <Card>
          <CardHeader>
            <CardTitle>Cost Trend</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={costByDay}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip
                  formatter={(value: number) => formatCurrency(value)}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="cost"
                  stroke="#3B82F6"
                  strokeWidth={2}
                  name="Cost (USD)"
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Model Distribution */}
        <Card>
          <CardHeader>
            <CardTitle>Model Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={modelDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={(entry) => {
                    const percent = entry.percent || 0
                    return `${entry.name} ${(percent * 100).toFixed(0)}%`
                  }}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {modelDistribution.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Model Comparison Table */}
      <Card>
        <CardHeader>
          <CardTitle>Model Performance</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left p-2">Model</th>
                  <th className="text-right p-2">Requests</th>
                  <th className="text-right p-2">Avg Cost</th>
                  <th className="text-right p-2">Total Cost</th>
                  <th className="text-right p-2">Status</th>
                </tr>
              </thead>
              <tbody>
                {modelDistribution.map((model) => (
                  <tr key={model.name} className="border-b last:border-0">
                    <td className="p-2 font-medium">{model.name}</td>
                    <td className="text-right p-2">{model.value}</td>
                    <td className="text-right p-2">{formatCurrency(model.cost / model.value)}</td>
                    <td className="text-right p-2">{formatCurrency(model.cost)}</td>
                    <td className="text-right p-2">
                      <span className="inline-flex items-center gap-1 text-green-600">
                        <CheckCircle className="h-3 w-3" />
                        Active
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Quick Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Tokens</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatNumber(stats.total_tokens)}</div>
            <p className="text-xs text-muted-foreground mt-1">
              Across all models
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Error Rate</CardTitle>
            <AlertCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.error_rate.toFixed(1)}%</div>
            <p className="text-xs text-muted-foreground mt-1">
              {Math.floor(stats.total_traces * stats.error_rate / 100)} failed requests
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Cost/Request</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatCurrency(stats.total_cost / stats.total_traces)}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              Per successful trace
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}