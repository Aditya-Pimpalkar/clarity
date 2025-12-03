import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'

export default function TracesPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Traces</h1>
        <p className="text-muted-foreground mt-2">
          View and debug your LLM traces
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Trace Explorer</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            🔍 Full trace exploration with filtering and search coming in Day 5!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
