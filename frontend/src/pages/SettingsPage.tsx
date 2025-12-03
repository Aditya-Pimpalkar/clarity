import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-2">
          Manage your account and preferences
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Settings & Configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            ⚙️ Settings, API keys, and preferences coming in Day 6!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
