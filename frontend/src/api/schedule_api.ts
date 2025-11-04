// API pour la gestion des schedules - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { Schedule, ScheduleCreateRequest, ScheduleUpdateRequest, ScheduleResponse, ScheduleListResponse } from '@/types/schedule'

/**
 * Récupère tous les schedules de l'utilisateur
 */
export async function getSchedules(): Promise<Schedule[]> {
  const { data } = await apiClient.get<ScheduleListResponse>('/api/schedules')
  return data.schedules || []
}

/**
 * Récupère un schedule par son ID
 */
export async function getScheduleById(id: number): Promise<Schedule> {
  const { data } = await apiClient.get<ScheduleResponse>(`/api/schedules/${id}`)
  return data.schedule
}

/**
 * Crée un nouveau schedule
 */
export async function createSchedule(scheduleData: ScheduleCreateRequest): Promise<Schedule> {
  const { data } = await apiClient.post<ScheduleResponse>('/api/schedules', scheduleData)
  return data.schedule
}

/**
 * Met à jour un schedule
 */
export async function updateSchedule(id: number, scheduleData: ScheduleUpdateRequest): Promise<Schedule> {
  const { data } = await apiClient.put<ScheduleResponse>(`/api/schedules/${id}`, scheduleData)
  return data.schedule
}

/**
 * Supprime un schedule
 */
export async function deleteSchedule(id: number): Promise<void> {
  await apiClient.delete(`/api/schedules/${id}`)
}