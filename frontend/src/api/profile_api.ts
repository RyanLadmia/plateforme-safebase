// API calls for user profile management
import { apiClient } from './axios'
import type { User } from '@/types/user'

export interface UpdateProfileRequest {
  firstname: string
  lastname: string
  email: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
  confirm_password: string
}

export interface ProfileResponse {
  message?: string
  user: User
}

/**
 * Get current user profile
 */
export async function getProfile(): Promise<User> {
  const { data } = await apiClient.get<{ user: User }>('/api/profile')
  return data.user
}

/**
 * Update user profile
 */
export async function updateProfile(profileData: UpdateProfileRequest): Promise<User> {
  const { data } = await apiClient.put<ProfileResponse>('/api/profile', profileData)
  return data.user
}

/**
 * Change user password
 */
export async function changePassword(passwordData: ChangePasswordRequest): Promise<void> {
  await apiClient.put('/api/profile/password', passwordData)
}

