import { create } from 'zustand'

interface AppState {
  sidebarOpen: boolean
  darkMode: boolean
  selectedOrg: string | null
  selectedProject: string | null
  setSidebarOpen: (open: boolean) => void
  toggleDarkMode: () => void
  setSelectedOrg: (orgId: string) => void
  setSelectedProject: (projectId: string) => void
}

export const useAppStore = create<AppState>((set) => ({
  sidebarOpen: true,
  darkMode: false,
  selectedOrg: null,
  selectedProject: null,

  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  
  toggleDarkMode: () => set((state) => {
    const newMode = !state.darkMode
    document.documentElement.classList.toggle('dark', newMode)
    return { darkMode: newMode }
  }),
  
  setSelectedOrg: (orgId) => set({ selectedOrg: orgId }),
  
  setSelectedProject: (projectId) => set({ selectedProject: projectId }),
}))
