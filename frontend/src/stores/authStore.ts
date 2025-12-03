import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '../types'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  checkAuth: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (email: string, _password: string) => {
        set({ isLoading: true, error: null })
        
        // MOCK LOGIN - Works without backend
        return new Promise((resolve) => {
          setTimeout(() => {
            const mockUser: User = {
              id: 'user-test-123',
              email: email,
              name: 'Diksha Sahare',
              organization_id: 'org-test-123',
              role: 'admin',
              created_at: new Date().toISOString(),
            }
            
            const mockToken = 'mock-jwt-token-123'
            localStorage.setItem('auth_token', mockToken)
            
            set({
              user: mockUser,
              token: mockToken,
              isAuthenticated: true,
              isLoading: false,
            })
            
            resolve()
          }, 300)
        })
      },

      logout: () => {
        localStorage.removeItem('auth_token')
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        })
      },

      checkAuth: async () => {
        const token = localStorage.getItem('auth_token')
        
        if (!token) {
          set({ isAuthenticated: false })
          return
        }
        
        // MOCK: Restore user from token without calling API
        const mockUser: User = {
          id: 'user-test-123',
          email: 'diksha@example.com',
          name: 'Diksha Sahare',
          organization_id: 'org-test-123',
          role: 'admin',
          created_at: new Date().toISOString(),
        }
        
        set({
          user: mockUser,
          token,
          isAuthenticated: true,
        })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        user: state.user,
      }),
    }
  )
)