// Types pour l'historique des actions utilisateur
export type ActivityType = 'all' | 'database' | 'backup' | 'schedule' | 'restore'

export interface HistoryItem {
  id: number
  action: string
  resource_type: 'database' | 'backup' | 'schedule' | 'restore'
  resource_id: number
  description: string
  created_at: string
  user_id: number
  metadata?: Record<string, any>
  ip_address?: string
  user_agent?: string
}

export interface HistoryResponse {
  history: HistoryItem[]
  total: number
  page: number
  limit: number
}

export interface HistoryFilters {
  type: ActivityType
  page: number
  limit: number
  start_date?: string
  end_date?: string
}

// Types pour les descriptions d'actions
export interface ActionDescription {
  icon: string
  color: string
  text: string
}

// Configuration des actions par type de ressource
export const ACTION_CONFIGS: Record<string, Record<string, ActionDescription>> = {
  database: {
    created: {
      icon: 'plus-circle',
      color: 'text-green-600',
      text: 'Base de données ajoutée'
    },
    updated: {
      icon: 'edit',
      color: 'text-blue-600',
      text: 'Base de données modifiée'
    },
    deleted: {
      icon: 'trash-2',
      color: 'text-red-600',
      text: 'Base de données supprimée'
    }
  },
  backup: {
    created: {
      icon: 'save',
      color: 'text-green-600',
      text: 'Sauvegarde créée'
    },
    completed: {
      icon: 'check-circle',
      color: 'text-green-600',
      text: 'Sauvegarde terminée'
    },
    failed: {
      icon: 'x-circle',
      color: 'text-red-600',
      text: 'Sauvegarde échouée'
    },
    deleted: {
      icon: 'trash-2',
      color: 'text-red-600',
      text: 'Sauvegarde supprimée'
    }
  },
  schedule: {
    created: {
      icon: 'clock',
      color: 'text-green-600',
      text: 'Planification créée'
    },
    updated: {
      icon: 'edit',
      color: 'text-blue-600',
      text: 'Planification modifiée'
    },
    deleted: {
      icon: 'trash-2',
      color: 'text-red-600',
      text: 'Planification supprimée'
    },
    executed: {
      icon: 'play',
      color: 'text-blue-600',
      text: 'Planification exécutée'
    }
  },
  restore: {
    created: {
      icon: 'rotate-ccw',
      color: 'text-green-600',
      text: 'Restauration effectuée'
    },
    restored: {
      icon: 'rotate-ccw',
      color: 'text-green-600',
      text: 'Restauration effectuée'
    },
    completed: {
      icon: 'check-circle',
      color: 'text-green-600',
      text: 'Restauration terminée'
    },
    failed: {
      icon: 'x-circle',
      color: 'text-red-600',
      text: 'Restauration échouée'
    },
    deleted: {
      icon: 'trash-2',
      color: 'text-red-600',
      text: 'Restauration supprimée'
    }
  }
}