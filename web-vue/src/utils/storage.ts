// Parse JWT token to get user_id
function parseJwt(token: string): { user_id?: number; exp?: number } {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
    return JSON.parse(jsonPayload)
  } catch (e) {
    return {}
  }
}

export const storage = {
  getToken: () => localStorage.getItem('token'),
  setToken: (token: string) => localStorage.setItem('token', token),
  removeToken: () => localStorage.removeItem('token'),

  getUserId: () => {
    const userId = localStorage.getItem('userId')
    if (userId) return userId

    // Try to get from token
    const token = localStorage.getItem('token')
    if (token) {
      const parsed = parseJwt(token)
      if (parsed.user_id) {
        return parsed.user_id.toString()
      }
    }
    return null
  },

  setUserId: (userId: string) => localStorage.setItem('userId', userId),
  removeUserId: () => localStorage.removeItem('userId'),

  getUserInfo: () => {
    const data = localStorage.getItem('userInfo')
    return data ? JSON.parse(data) : null
  },
  setUserInfo: (info: any) => localStorage.setItem('userInfo', JSON.stringify(info)),
  removeUserInfo: () => localStorage.removeItem('userInfo'),

  clear: () => localStorage.clear()
}
