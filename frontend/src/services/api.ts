import axios, { AxiosInstance, AxiosError } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('auth_token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        // DON'T auto-logout on 401 - let components handle it
        // This allows mock mode to work
        return Promise.reject(error)
      }
    )
  }

  async login(email: string, password: string) {
    const response = await this.client.post('/api/v1/auth/login', {
      email,
      password,
    })
    return response.data
  }

  async getCurrentUser() {
    const response = await this.client.get('/api/v1/auth/me')
    return response.data
  }

  async getDashboard(timeRange: string = '24h') {
    try {
      const response = await this.client.get('/api/v1/analytics/dashboard', {
        params: { time_range: timeRange },
      })
      return response.data
    } catch (error) {
      // Return mock data if API call fails
      return {
        data: {
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
        }
      }
    }
  }

  async getTraces(params: any) {
    const response = await this.client.get('/api/v1/traces', { params })
    return response.data
  }

  async getTrace(traceId: string) {
    const response = await this.client.get(`/api/v1/traces/${traceId}`)
    return response.data
  }

  async getCostAnalysis(timeRange: string = '30d') {
    const response = await this.client.get('/api/v1/analytics/costs', {
      params: { time_range: timeRange },
    })
    return response.data
  }

  async getPerformanceMetrics(timeRange: string = '24h') {
    const response = await this.client.get('/api/v1/analytics/performance', {
      params: { time_range: timeRange },
    })
    return response.data
  }

  async getModelComparison(timeRange: string = '7d') {
    const response = await this.client.get('/api/v1/analytics/models', {
      params: { time_range: timeRange },
    })
    return response.data
  }

  async checkHealth() {
    const response = await this.client.get('/health')
    return response.data
  }
}

export const api = new ApiClient()
export default api