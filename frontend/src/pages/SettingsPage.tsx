import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { 
  User,
  Key,
  Bell,
  Download,
  Trash2,
  Copy,
  Plus,
  Eye,
  EyeOff
} from 'lucide-react'
import { useAuthStore } from '../stores/authStore'
import { useAppStore } from '../stores/appStore'

export default function SettingsPage() {
  const user = useAuthStore((state) => state.user)
  const { darkMode, toggleDarkMode } = useAppStore()
  const [apiKeys, setApiKeys] = useState([
    { id: '1', name: 'Production Key', key: 'llm_prod_***************', created: '2025-12-01', lastUsed: '2 hours ago' },
    { id: '2', name: 'Development Key', key: 'llm_dev_***************', created: '2025-11-28', lastUsed: '1 day ago' },
  ])
  const [showKey, setShowKey] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const generateNewKey = () => {
    const newKey = {
      id: Date.now().toString(),
      name: 'New API Key',
      key: 'llm_new_' + Math.random().toString(36).substring(2, 15),
      created: new Date().toISOString().split('T')[0],
      lastUsed: 'Never',
    }
    setApiKeys([...apiKeys, newKey])
  }

  const deleteKey = (id: string) => {
    if (confirm('Are you sure you want to delete this API key?')) {
      setApiKeys(apiKeys.filter(k => k.id !== id))
    }
  }

  return (
    <div className="space-y-6 max-w-4xl">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-2">
          Manage your account and preferences
        </p>
      </div>

      {/* Profile Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="h-5 w-5" />
            Profile Information
          </CardTitle>
          <CardDescription>
            Update your personal information and preferences
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="text-sm font-medium">Name</label>
              <input
                type="text"
                value={user?.name || 'Diksha Sahare'}
                className="w-full mt-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
                readOnly
              />
            </div>
            <div>
              <label className="text-sm font-medium">Email</label>
              <input
                type="email"
                value={user?.email || 'diksha@example.com'}
                className="w-full mt-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
                readOnly
              />
            </div>
            <div>
              <label className="text-sm font-medium">Organization</label>
              <input
                type="text"
                value={user?.organization_id || 'org-test-123'}
                className="w-full mt-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
                readOnly
              />
            </div>
            <div>
              <label className="text-sm font-medium">Role</label>
              <input
                type="text"
                value={user?.role || 'admin'}
                className="w-full mt-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background capitalize"
                readOnly
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* API Keys */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Key className="h-5 w-5" />
                API Keys
              </CardTitle>
              <CardDescription>
                Manage API keys for programmatic access
              </CardDescription>
            </div>
            <Button onClick={generateNewKey} size="sm">
              <Plus className="h-4 w-4 mr-2" />
              Generate New Key
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {apiKeys.map((apiKey) => (
              <div key={apiKey.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h4 className="font-medium">{apiKey.name}</h4>
                    <span className="text-xs text-muted-foreground">
                      Created {apiKey.created}
                    </span>
                  </div>
                  <div className="flex items-center gap-2 mt-2">
                    <code className="text-sm font-mono bg-muted px-2 py-1 rounded">
                      {showKey === apiKey.id ? apiKey.key : apiKey.key}
                    </code>
                    <button
                      onClick={() => copyToClipboard(apiKey.key)}
                      className="p-1 hover:bg-accent rounded"
                    >
                      <Copy className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => setShowKey(showKey === apiKey.id ? null : apiKey.id)}
                      className="p-1 hover:bg-accent rounded"
                    >
                      {showKey === apiKey.id ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </button>
                  </div>
                  <p className="text-xs text-muted-foreground mt-2">
                    Last used: {apiKey.lastUsed}
                  </p>
                </div>
                <button
                  onClick={() => deleteKey(apiKey.id)}
                  className="p-2 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-md text-red-600"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            ))}
          </div>
          
          {copied && (
            <div className="mt-3 p-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded text-sm text-green-600 text-center">
              ✓ Copied to clipboard!
            </div>
          )}
        </CardContent>
      </Card>

      {/* Appearance */}
      <Card>
        <CardHeader>
          <CardTitle>Appearance</CardTitle>
          <CardDescription>
            Customize how the application looks
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium">Dark Mode</h4>
              <p className="text-sm text-muted-foreground mt-1">
                Toggle dark mode for better visibility
              </p>
            </div>
            <button
              onClick={toggleDarkMode}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                darkMode ? 'bg-primary' : 'bg-gray-300'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  darkMode ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>
        </CardContent>
      </Card>

      {/* Notifications */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Bell className="h-5 w-5" />
            Notifications
          </CardTitle>
          <CardDescription>
            Configure alert preferences
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium">Error Alerts</h4>
              <p className="text-sm text-muted-foreground mt-1">
                Get notified when error rate exceeds 5%
              </p>
            </div>
            <button className="relative inline-flex h-6 w-11 items-center rounded-full bg-primary">
              <span className="inline-block h-4 w-4 transform rounded-full bg-white translate-x-6" />
            </button>
          </div>

          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium">Cost Alerts</h4>
              <p className="text-sm text-muted-foreground mt-1">
                Alert when daily cost exceeds $0.50
              </p>
            </div>
            <button className="relative inline-flex h-6 w-11 items-center rounded-full bg-gray-300">
              <span className="inline-block h-4 w-4 transform rounded-full bg-white translate-x-1" />
            </button>
          </div>

          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium">Performance Alerts</h4>
              <p className="text-sm text-muted-foreground mt-1">
                Notify when P95 latency exceeds 2000ms
              </p>
            </div>
            <button className="relative inline-flex h-6 w-11 items-center rounded-full bg-primary">
              <span className="inline-block h-4 w-4 transform rounded-full bg-white translate-x-6" />
            </button>
          </div>
        </CardContent>
      </Card>

      {/* Data Management */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Download className="h-5 w-5" />
            Data Management
          </CardTitle>
          <CardDescription>
            Export or manage your data
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Button variant="outline" className="w-full justify-start">
            <Download className="h-4 w-4 mr-2" />
            Export Traces (CSV)
          </Button>
          <Button variant="outline" className="w-full justify-start">
            <Download className="h-4 w-4 mr-2" />
            Export Analytics Report (PDF)
          </Button>
          <Button variant="outline" className="w-full justify-start text-red-600 hover:text-red-600">
            <Trash2 className="h-4 w-4 mr-2" />
            Clear Old Data (90+ days)
          </Button>
        </CardContent>
      </Card>

      {/* Danger Zone */}
      <Card className="border-red-200 dark:border-red-800">
        <CardHeader>
          <CardTitle className="text-red-600">Danger Zone</CardTitle>
          <CardDescription>
            Irreversible actions - proceed with caution
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button variant="destructive" className="w-full">
            <Trash2 className="h-4 w-4 mr-2" />
            Delete All Data
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}