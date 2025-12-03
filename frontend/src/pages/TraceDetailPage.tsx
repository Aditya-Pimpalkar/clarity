import { useParams } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'

export default function TraceDetailPage() {
  const { id } = useParams()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Trace Details</h1>
        <p className="text-muted-foreground mt-2">
          Trace ID: {id}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Detailed View</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            🔬 Detailed trace view with span timeline coming in Day 5!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
