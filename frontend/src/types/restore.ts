// Types pour la gestion des restaurations
export interface Restore {
  id: number
  status: 'pending' | 'running' | 'success' | 'failed'
  created_at: string
  updated_at: string
  user_id: number
  backup_id: number
  database_id: number
  user?: {
    id: number
    firstname: string
    lastname: string
    email: string
  }
  backup?: {
    id: number
    filename: string
    status: string
    created_at: string
  }
  database?: {
    id: number
    name: string
    type: string
    host: string
    port: string
  }
}

export interface RestoreResponse {
  restore: Restore
}

export interface RestoreListResponse {
  restores: Restore[]
}