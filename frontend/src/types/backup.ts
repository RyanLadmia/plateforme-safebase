// Types pour les sauvegardes
export interface Backup {
  id: number
  filename: string
  filepath: string
  size: number
  status: 'pending' | 'completed' | 'failed'
  error_msg?: string
  created_at: string
  updated_at: string
  user_id: number
  database_id: number
  user_agent?: string
}

export interface BackupResponse {
  backup: Backup
  message?: string
}

export interface BackupListResponse {
  backups: Backup[]
}
