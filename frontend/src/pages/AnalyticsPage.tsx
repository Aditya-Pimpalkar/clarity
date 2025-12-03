import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'

export default function AnalyticsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Analytics</h1>
        <p className="text-muted-foreground mt-2">
          Deep dive into your LLM metrics
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Advanced Analytics</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            📈 Cost analysis, performance metrics, and model comparisons coming in Day 4!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
