// API pour la gestion des bases de données - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { 
  Database, 
  DatabaseCreateRequest, 
  DatabaseUpdateRequest,
  DatabaseResponse,
  DatabaseListResponse 
} from '@/types/database'

/**
 * Récupère toutes les bases de données de l'utilisateur
 */
export async function getDatabases(): Promise<Database[]> {
  const { data } = await apiClient.get<DatabaseListResponse>('/api/databases')
  return data.databases || []
}

/**
 * Récupère une base de données par son ID
 */
export async function getDatabaseById(id: number): Promise<Database> {
  const { data } = await apiClient.get<DatabaseResponse>(`/api/databases/${id}`)
  return data.database
}

/**
 * Crée une nouvelle base de données
 */
export async function createDatabase(databaseData: DatabaseCreateRequest): Promise<Database> {
  const { data } = await apiClient.post<DatabaseResponse>('/api/databases', databaseData)
  return data.database
}

/**
 * Met à jour une base de données
 */
export async function updateDatabase(id: number, databaseData: DatabaseUpdateRequest): Promise<Database> {
  const { data } = await apiClient.put<DatabaseResponse>(`/api/databases/${id}`, databaseData)
  return data.database
}

/**
 * Supprime une base de données
 */
export async function deleteDatabase(id: number): Promise<void> {
  await apiClient.delete(`/api/databases/${id}`)
}
