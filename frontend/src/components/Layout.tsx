import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import Header from './Header'
import { useAppStore } from '../stores/appStore'

export default function Layout() {
  const sidebarOpen = useAppStore((state) => state.sidebarOpen)

  return (
    <div className="min-h-screen bg-background">
      <Sidebar />
      
      <div className={`transition-all duration-300 ${sidebarOpen ? 'lg:pl-64' : 'lg:pl-20'}`}>
        <Header />
        
        <main className="p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
