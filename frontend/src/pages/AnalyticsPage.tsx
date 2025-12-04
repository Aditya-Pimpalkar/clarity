import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { 
  DollarSign,
  TrendingUp,
  TrendingDown,
  Zap,
  Clock,
  BarChart3,
} from 'lucide-react'
import {
  BarChart,
  Bar,
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
  RadarChart,
  Radar,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
} from 'recharts'
import { formatCurrency, formatDuration, formatNumber } from '../lib/utils'

export default function AnalyticsPage() {
  const [timeRange, setTimeRange] = useState('7d')

  // Mock data - will be replaced with real API calls
  const costByModel = [
    { model: 'GPT-4', cost: 0.089, requests: 45, avgCost: 0.00198 },
    { model: 'GPT-3.5', cost: 0.002, requests: 12, avgCost: 0.00017 },
    { model: 'Claude-3-Sonnet', cost: 0.059, requests: 7, avgCost: 0.00843 },
  ]

  const costByProvider = [
    { name: 'OpenAI', value: 0.091, percentage: 61 },
    { name: 'Anthropic', value: 0.059, percentage: 39 },
  ]

  const performanceTrend = [
    { date: 'Mon', latency: 145, requests: 12, errors: 0 },
    { date: 'Tue', latency: 132, requests: 18, errors: 1 },
    { date: 'Wed', latency: 156, requests: 15, errors: 0 },
    { date: 'Thu', latency: 128, requests: 21, errors: 0 },
    { date: 'Fri', latency: 149, requests: 19, errors: 2 },
    { date: 'Sat', latency: 138, requests: 8, errors: 0 },
    { date: 'Sun', latency: 142, requests: 6, errors: 0 },
  ]

  const modelEfficiency = [
    { model: 'GPT-4', performance: 85, cost: 65, reliability: 95, speed: 70 },
    { model: 'GPT-3.5', performance: 70, cost: 95, reliability: 90, speed: 85 },
    { model: 'Claude-3', performance: 90, cost: 70, reliability: 98, speed: 65 },
  ]

  const tokenUsage = [
    { date: 'Mon', prompt: 1200, completion: 800 },
    { date: 'Tue', prompt: 1500, completion: 950 },
    { date: 'Wed', prompt: 1300, completion: 850 },
    { date: 'Thu', prompt: 1800, completion: 1100 },
    { date: 'Fri', prompt: 1600, completion: 1000 },
    { date: 'Sat', prompt: 800, completion: 500 },
    { date: 'Sun', prompt: 600, completion: 400 },
  ]

  const COLORS = {
    primary: '#3B82F6',
    success: '#10B981',
    warning: '#F59E0B',
    danger: '#EF4444',
    purple: '#8B5CF6',
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Analytics</h1>
          <p className="text-muted-foreground mt-2">
            Deep dive into your LLM metrics and costs
          </p>
        </div>

        {/* Time Range Selector */}
        <div className="flex gap-2">
          {['24h', '7d', '30d', '90d'].map((range) => (
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

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Spend</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$0.15</div>
            <div className="flex items-center gap-1 text-xs text-green-600 mt-1">
              <TrendingDown className="h-3 w-3" />
              <span>12% vs last period</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Cost/Request</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$0.0023</div>
            <p className="text-xs text-muted-foreground mt-1">
              Across all models
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Monthly Projection</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$4.50</div>
            <p className="text-xs text-muted-foreground mt-1">
              Based on current usage
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cost Efficiency</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">87%</div>
            <div className="flex items-center gap-1 text-xs text-green-600 mt-1">
              <TrendingUp className="h-3 w-3" />
              <span>Excellent</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Cost Analysis */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Cost by Model */}
        <Card>
          <CardHeader>
            <CardTitle>Cost Breakdown by Model</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={costByModel}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="model" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip
                  formatter={(value: number, name: string) => {
                    if (name === 'cost') return [formatCurrency(value), 'Cost']
                    if (name === 'avgCost') return [formatCurrency(value), 'Avg Cost']
                    return [value, name]
                  }}
                />
                <Legend />
                <Bar dataKey="cost" fill={COLORS.primary} name="Total Cost" />
                <Bar dataKey="avgCost" fill={COLORS.success} name="Avg Cost" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Cost by Provider */}
        <Card>
          <CardHeader>
            <CardTitle>Cost Distribution by Provider</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={costByProvider}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={(entry) => `${entry.name} ${((entry.percent || 0) * 100).toFixed(0)}%`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {costByProvider.map((_, index) => (
                    <Cell 
                      key={`cell-${index}`} 
                      fill={index === 0 ? COLORS.primary : COLORS.success} 
                    />
                  ))}
                </Pie>
                <Tooltip formatter={(value: number) => formatCurrency(value)} />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Performance Trends */}
      <Card>
        <CardHeader>
          <CardTitle>Performance & Usage Trends</CardTitle>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={performanceTrend}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="latency"
                stroke={COLORS.primary}
                strokeWidth={2}
                name="Latency (ms)"
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="requests"
                stroke={COLORS.success}
                strokeWidth={2}
                name="Requests"
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="errors"
                stroke={COLORS.danger}
                strokeWidth={2}
                name="Errors"
              />
            </LineChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Model Efficiency Radar */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Model Efficiency Comparison</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <RadarChart data={modelEfficiency}>
                <PolarGrid />
                <PolarAngleAxis dataKey="model" />
                <PolarRadiusAxis angle={90} domain={[0, 100]} />
                <Radar
                  name="Performance"
                  dataKey="performance"
                  stroke={COLORS.primary}
                  fill={COLORS.primary}
                  fillOpacity={0.3}
                />
                <Radar
                  name="Cost Efficiency"
                  dataKey="cost"
                  stroke={COLORS.success}
                  fill={COLORS.success}
                  fillOpacity={0.3}
                />
                <Radar
                  name="Reliability"
                  dataKey="reliability"
                  stroke={COLORS.purple}
                  fill={COLORS.purple}
                  fillOpacity={0.3}
                />
                <Legend />
                <Tooltip />
              </RadarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Token Usage */}
        <Card>
          <CardHeader>
            <CardTitle>Token Usage Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={tokenUsage}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="prompt" stackId="a" fill={COLORS.primary} name="Prompt Tokens" />
                <Bar dataKey="completion" stackId="a" fill={COLORS.success} name="Completion Tokens" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Model Performance Table */}
      <Card>
        <CardHeader>
          <CardTitle>Detailed Model Metrics</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left p-3 text-sm font-medium text-muted-foreground">Model</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Requests</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Total Cost</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Avg Cost</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Avg Latency</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Success Rate</th>
                  <th className="text-right p-3 text-sm font-medium text-muted-foreground">Efficiency</th>
                </tr>
              </thead>
              <tbody>
                {costByModel.map((model) => (
                  <tr key={model.model} className="border-b last:border-0 hover:bg-muted/50">
                    <td className="p-3 font-medium">{model.model}</td>
                    <td className="text-right p-3">{formatNumber(model.requests)}</td>
                    <td className="text-right p-3">{formatCurrency(model.cost)}</td>
                    <td className="text-right p-3">{formatCurrency(model.avgCost)}</td>
                    <td className="text-right p-3">
                      <div className="flex items-center justify-end gap-1">
                        <Clock className="h-3 w-3 text-muted-foreground" />
                        {formatDuration(145)}
                      </div>
                    </td>
                    <td className="text-right p-3">
                      <span className="text-green-600 font-medium">98.5%</span>
                    </td>
                    <td className="text-right p-3">
                      <div className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 text-xs font-medium">
                        <Zap className="h-3 w-3" />
                        High
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Cost Insights */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">💡 Cost Optimization</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Switching 30% of GPT-4 requests to GPT-3.5 could save{' '}
              <span className="font-bold text-green-600">$0.027/day</span>
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">⚡ Performance Insight</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Your average latency is{' '}
              <span className="font-bold text-green-600">23% faster</span>{' '}
              than industry average
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">🎯 Usage Pattern</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Peak usage: <span className="font-bold">Thu-Fri</span>.
              Consider rate limiting or caching strategies.
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Cost Projection */}
      <Card>
        <CardHeader>
          <CardTitle>Cost Projection</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-muted rounded-lg">
              <div>
                <p className="text-sm text-muted-foreground">Daily Average</p>
                <p className="text-2xl font-bold">{formatCurrency(0.021)}</p>
              </div>
              <div className="text-right">
                <p className="text-sm text-muted-foreground">Weekly</p>
                <p className="text-xl font-bold">{formatCurrency(0.15)}</p>
              </div>
              <div className="text-right">
                <p className="text-sm text-muted-foreground">Monthly</p>
                <p className="text-xl font-bold">{formatCurrency(0.64)}</p>
              </div>
              <div className="text-right">
                <p className="text-sm text-muted-foreground">Yearly</p>
                <p className="text-xl font-bold">{formatCurrency(7.68)}</p>
              </div>
            </div>

            <div className="grid gap-3 md:grid-cols-3">
              <div className="p-3 border rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-2 h-2 rounded-full bg-green-500"></div>
                  <span className="text-sm font-medium">Best Case</span>
                </div>
                <p className="text-lg font-bold text-green-600">{formatCurrency(5.12)}/year</p>
                <p className="text-xs text-muted-foreground mt-1">With optimization</p>
              </div>

              <div className="p-3 border rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-2 h-2 rounded-full bg-blue-500"></div>
                  <span className="text-sm font-medium">Expected</span>
                </div>
                <p className="text-lg font-bold text-blue-600">{formatCurrency(7.68)}/year</p>
                <p className="text-xs text-muted-foreground mt-1">Current trajectory</p>
              </div>

              <div className="p-3 border rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-2 h-2 rounded-full bg-orange-500"></div>
                  <span className="text-sm font-medium">Worst Case</span>
                </div>
                <p className="text-lg font-bold text-orange-600">{formatCurrency(12.24)}/year</p>
                <p className="text-xs text-muted-foreground mt-1">If usage grows 60%</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Recommendations */}
      <Card>
        <CardHeader>
          <CardTitle>💡 Recommendations</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex items-start gap-3 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
              <div className="w-6 h-6 rounded-full bg-green-600 flex items-center justify-center flex-shrink-0 mt-0.5">
                <span className="text-white text-xs font-bold">1</span>
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-sm">Implement Caching</h4>
                <p className="text-sm text-muted-foreground mt-1">
                  38% of requests have similar inputs. Caching could save ~$0.03/day
                </p>
              </div>
              <span className="text-xs text-green-600 font-medium">High Impact</span>
            </div>

            <div className="flex items-start gap-3 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
              <div className="w-6 h-6 rounded-full bg-blue-600 flex items-center justify-center flex-shrink-0 mt-0.5">
                <span className="text-white text-xs font-bold">2</span>
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-sm">Optimize GPT-4 Usage</h4>
                <p className="text-sm text-muted-foreground mt-1">
                  Use GPT-3.5 for simple queries. Could reduce costs by 40%
                </p>
              </div>
              <span className="text-xs text-blue-600 font-medium">Medium Impact</span>
            </div>

            <div className="flex items-start gap-3 p-3 bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg">
              <div className="w-6 h-6 rounded-full bg-purple-600 flex items-center justify-center flex-shrink-0 mt-0.5">
                <span className="text-white text-xs font-bold">3</span>
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-sm">Prompt Engineering</h4>
                <p className="text-sm text-muted-foreground mt-1">
                  Reduce average prompt tokens by 15% with better prompt design
                </p>
              </div>
              <span className="text-xs text-purple-600 font-medium">Low Impact</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}