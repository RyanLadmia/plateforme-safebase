// API calls for action history
import axios from '@/api/axios'
import type { HistoryResponse, HistoryFilters } from '@/types/history'

/**
 * Get user action history
 */
export async function getUserActionHistory(page: number = 1, limit: number = 20): Promise<HistoryResponse> {
  const response = await axios.get('/api/history', {
    params: { page, limit }
  })
  return response.data
}

/**
 * Get action history by resource type
 */
export async function getActionHistoryByType(resourceType: string, page: number = 1, limit: number = 20): Promise<HistoryResponse> {
  const response = await axios.get(`/api/history/type/${resourceType}`, {
    params: { page, limit }
  })
  return response.data
}

/**
 * Get action history for a specific resource
 */
export async function getResourceActionHistory(resourceType: string, resourceId: number): Promise<{ history: any[] }> {
  const response = await axios.get(`/api/history/resource/${resourceType}/${resourceId}`)
  return response.data
}

/**
 * Get recent action history (admin only)
 */
export async function getRecentActionHistory(page: number = 1, limit: number = 20): Promise<HistoryResponse> {
  const response = await axios.get('/api/history/recent', {
    params: { page, limit }
  })
  return response.data
}