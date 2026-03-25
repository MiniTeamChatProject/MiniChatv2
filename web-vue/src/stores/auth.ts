import { reactive, computed } from 'vue'
import { userApi } from '@/api/user'
import { storage } from '@/utils/storage'

// Parse JWT token to get user_id
function parseJwt(token: string): { user_id?: number; exp?: number } | null {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
    return JSON.parse(jsonPayload)
  } catch (e) {
    return null
  }
}

// Global auth state - single instance
const authState = reactive({
  token: storage.getToken() || '',
  userId: storage.getUserId() || '',
  userInfo: storage.getUserInfo() || null,
  _initialized: false
})

// Initialize auth state from storage on module load
async function initializeAuthState() {
  if (authState._initialized) return

  const token = storage.getToken()
  const userId = storage.getUserId()
  const userInfo = storage.getUserInfo()

  if (token) {
    authState.token = token
    if (userId) {
      authState.userId = userId
    }
    if (userInfo) {
      authState.userInfo = userInfo
    }
    // Verify token is still valid by fetching user info
    try {
      const info = await userApi.getProfile()
      authState.userInfo = info
      storage.setUserInfo(info)
    } catch (error) {
      // Token might be expired, clear auth state
      console.warn('Failed to verify token, clearing auth state:', error)
      clearAuthState()
    }
  }
  authState._initialized = true
}

function clearAuthState() {
  authState.token = ''
  authState.userId = ''
  authState.userInfo = null
  storage.clear()
}

// Export the auth store as a singleton
export const authStore = {
  state: authState,

  get isLoggedIn(): boolean {
    return !!authState.token
  },

  get token(): string {
    return authState.token
  },

  get userId(): string {
    return authState.userId
  },

  get userInfo(): any {
    return authState.userInfo
  },

  get username(): string {
    return authState.userInfo?.username || ''
  },

  async login(username: string, password: string) {
    const response = await userApi.login({ username, password })
    authState.token = response.token
    // Extract user_id from JWT token
    const parsed = parseJwt(response.token)
    if (parsed?.user_id) {
      authState.userId = parsed.user_id.toString()
      storage.setUserId(parsed.user_id.toString())
    }
    storage.setToken(response.token)
    return response
  },

  async register(username: string, password: string, nickname: string) {
    const response = await userApi.register({ username, password, nickname })
    return response
  },

  async fetchUserInfo() {
    const info = await userApi.getProfile()
    authState.userInfo = info
    storage.setUserInfo(info)
    return info
  },

  logout() {
    clearAuthState()
  },

  persistAuth() {
    if (authState.token) storage.setToken(authState.token)
    if (authState.userId) storage.setUserId(authState.userId)
    if (authState.userInfo) storage.setUserInfo(authState.userInfo)
  },

  initialize: initializeAuthState
}

// Type export for use in components
export type AuthStore = typeof authStore
