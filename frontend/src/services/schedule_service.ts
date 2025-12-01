// Service de gestion des schedules - Logique métier
import * as scheduleApi from '@/api/schedule_api'
import type { Schedule, ScheduleCreateRequest, ScheduleUpdateRequest } from '@/types/schedule'

/**
 * Service de gestion des schedules
 */
export class ScheduleService {
  /**
   * Récupère tous les schedules
   */
  async fetchSchedules(): Promise<Schedule[]> {
    return await scheduleApi.getSchedules()
  }

  /**
   * Récupère un schedule par son ID
   */
  async fetchScheduleById(id: number): Promise<Schedule> {
    return await scheduleApi.getScheduleById(id)
  }

  /**
   * Crée un nouveau schedule
   */
  async createSchedule(scheduleData: ScheduleCreateRequest): Promise<Schedule> {
    if (!scheduleData.database_id || scheduleData.database_id <= 0) {
      throw new Error('ID de base de données invalide')
    }
    if (!scheduleData.cron_expression || scheduleData.cron_expression.trim() === '') {
      throw new Error('Expression cron requise')
    }

    return await scheduleApi.createSchedule(scheduleData)
  }

  /**
   * Met à jour un schedule
   */
  async updateSchedule(id: number, scheduleData: ScheduleUpdateRequest): Promise<Schedule> {
    if (!id || id <= 0) {
      throw new Error('ID de schedule invalide')
    }

    return await scheduleApi.updateSchedule(id, scheduleData)
  }

  /**
   * Active/désactive un schedule
   */
  async toggleSchedule(id: number, active: boolean): Promise<Schedule> {
    return await this.updateSchedule(id, { active })
  }

  /**
   * Met à jour l'expression cron d'un schedule
   */
  async updateCronExpression(id: number, cronExpression: string): Promise<Schedule> {
    if (!cronExpression || cronExpression.trim() === '') {
      throw new Error('Expression cron requise')
    }

    return await this.updateSchedule(id, { cron_expression: cronExpression })
  }

  /**
   * Supprime un schedule
   */
  async deleteSchedule(id: number): Promise<void> {
    if (!id || id <= 0) {
      throw new Error('ID de schedule invalide')
    }

    await scheduleApi.deleteSchedule(id)
  }

  /**
   * Valide une expression cron basique
   */
  validateCronExpression(expression: string): boolean {
    if (!expression || expression.trim() === '') {
      return false
    }

    // Validation basique : doit contenir 5 parties séparées par des espaces
    const parts = expression.trim().split(/\s+/)
    return parts.length === 5
  }
}

// Export d'une instance unique du service
export const scheduleService = new ScheduleService()