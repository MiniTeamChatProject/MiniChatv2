import { userHttp } from './http'
import type { LoginRequest, RegisterRequest, LoginResponse, UserProfile } from '@/types'

export const userApi = {
  // Login
  login: (data: LoginRequest) =>
    userHttp.post<LoginResponse>('/login', data),

  // Register
  register: (data: RegisterRequest) =>
    userHttp.post<{ msg: string }>('/register', data),

  // Get user profile
  getProfile: () =>
    userHttp.get<UserProfile>('/user/profile'),
}
