// Types pour les bases de données
export interface Database {
  id: number
  name: string
  type: 'mysql' | 'postgresql'
  host: string
  port: string
  username: string
  db_name: string
  created_at: string
  updated_at: string
  user_id: number
}

export interface DatabaseCreateRequest {
  name: string
  type: 'mysql' | 'postgresql'
  host: string
  port: string
  username: string
  password: string
  db_name: string
}

export interface DatabaseUpdateRequest {
  name: string
  type: 'mysql' | 'postgresql'
  host: string
  port: string
  username: string
  password?: string // Optionnel lors de la mise à jour
  db_name: string
}

export interface DatabaseResponse {
  database: Database
  message?: string
}

export interface DatabaseListResponse {
  databases: Database[]
}
