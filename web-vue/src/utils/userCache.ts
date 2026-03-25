import { authStore } from '@/stores/auth'

// Simple in-memory username cache
const usernameCache = new Map<number, string>()

export function getUsername(userId: number): string {
  // Check cache first
  if (usernameCache.has(userId)) {
    return usernameCache.get(userId)!
  }

  // Check if it's current user
  if (String(userId) === String(authStore.userId)) {
    // Use nickname first, then username as fallback
    const nickname = authStore.userInfo?.nickname
    const username = authStore.username || authStore.userInfo?.username
    const displayName = nickname || username || `User${userId}`
    if (displayName) {
      usernameCache.set(userId, displayName)
      return displayName
    }
  }

  // Default fallback
  return `User${userId}`
}

export function cacheUsername(userId: number, username: string) {
  usernameCache.set(userId, username)
}

// Pre-fetch usernames for a list of user IDs
export async function prefetchUsernames(userIds: number[]) {
  // Filter out users we already have cached
  const uncachedIds = userIds.filter(id => !usernameCache.has(id))

  if (uncachedIds.length === 0) return

  // For now, just use default usernames
  // In the future, this could call the user service
  for (const id of uncachedIds) {
    if (String(id) === String(authStore.userId)) {
      const nickname = authStore.userInfo?.nickname
      const username = authStore.username || authStore.userInfo?.username
      const displayName = nickname || username
      if (displayName) {
        usernameCache.set(id, displayName)
      }
    }
  }
}

export function clearUserCache() {
  usernameCache.clear()
}
