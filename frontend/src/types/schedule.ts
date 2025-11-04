// Types pour les schedules (tâches planifiées)
export interface Schedule {
  id: number
  cron_expression: string
  active: boolean
  created_at: string
  updated_at: string
  user_id: number
  database_id: number
  database?: {
    id: number
    name: string
    type: 'mysql' | 'postgresql'
  }
}

export interface ScheduleCreateRequest {
  database_id: number
  cron_expression: string
}

export interface ScheduleUpdateRequest {
  cron_expression?: string
  active?: boolean
}

export interface ScheduleResponse {
  schedule: Schedule
  message?: string
}

export interface ScheduleListResponse {
  schedules: Schedule[]
}

// Types pour les expressions cron prédéfinies
export interface CronPreset {
  label: string
  description: string
  expression: string
}

export const CRON_PRESETS: CronPreset[] = [
  {
    label: 'Toutes les heures',
    description: 'Exécution toutes les heures',
    expression: '0 * * * *'
  },
  {
    label: 'Tous les jours à minuit',
    description: 'Exécution quotidienne à 00:00',
    expression: '0 0 * * *'
  },
  {
    label: 'Tous les jours à 6h',
    description: 'Exécution quotidienne à 06:00',
    expression: '0 6 * * *'
  },
  {
    label: 'Tous les jours à midi',
    description: 'Exécution quotidienne à 12:00',
    expression: '0 12 * * *'
  },
  {
    label: 'Tous les jours à 18h',
    description: 'Exécution quotidienne à 18:00',
    expression: '0 18 * * *'
  },
  {
    label: 'Toutes les semaines (lundi)',
    description: 'Exécution hebdomadaire le lundi à minuit',
    expression: '0 0 * * 1'
  },
  {
    label: 'Tous les mois (1er)',
    description: 'Exécution mensuelle le 1er à minuit',
    expression: '0 0 1 * *'
  }
]