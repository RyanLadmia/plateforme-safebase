// Types pour les bases de données
export interface Database {
  id: number
  name: string
  type: 'mysql' | 'postgresql'
  host: string
  port: string
  username: string
  db_name: string
  url?: string // URL complète (optionnel)
  created_at: string
  updated_at: string
  user_id: number
}

export interface DatabaseCreateRequest {
  name: string
  type: 'mysql' | 'postgresql'
  host?: string // Optionnel si URL fournie
  port?: string // Optionnel si URL fournie
  username?: string // Optionnel si URL fournie
  password?: string // Optionnel si URL fournie
  db_name?: string // Optionnel si URL fournie
  url?: string // Alternative: URL complète
}

export interface DatabaseUpdateRequest {
  name: string
  type: 'mysql' | 'postgresql'
  host?: string // Optionnel si URL fournie
  port?: string // Optionnel si URL fournie
  username?: string // Optionnel si URL fournie
  password?: string // Optionnel lors de la mise à jour
  db_name?: string // Optionnel si URL fournie
  url?: string // Alternative: URL complète
}

export interface DatabaseResponse {
  database: Database
  message?: string
}

export interface DatabaseListResponse {
  databases: Database[]
}
