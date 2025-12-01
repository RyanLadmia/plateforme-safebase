// Utilitaires pour analyser et décrire les expressions cron
export class CronUtils {
  /**
   * Convertit une expression cron en description française compréhensible
   */
  static getFrequencyDescription(cronExpression: string): string {
    try {
      const parts = cronExpression.trim().split(/\s+/)
      if (parts.length !== 5) return 'Expression invalide'

      const [minute, hour, day, month, dayOfWeek] = parts

      // Expressions prédéfinies courantes
      const commonExpressions: { [key: string]: string } = {
        '0 0 * * *': 'Tous les jours à minuit',
        '0 6 * * *': 'Tous les jours à 6h',
        '0 12 * * *': 'Tous les jours à midi',
        '0 18 * * *': 'Tous les jours à 18h',
        '0 */6 * * *': 'Toutes les 6 heures',
        '0 */12 * * *': 'Toutes les 12 heures',
        '0 * * * *': 'Toutes les heures',
        '*/30 * * * *': 'Toutes les 30 minutes',
        '*/15 * * * *': 'Toutes les 15 minutes',
        '*/10 * * * *': 'Toutes les 10 minutes',
        '*/5 * * * *': 'Toutes les 5 minutes',
        '0 0 * * 1': 'Tous les lundis à minuit',
        '0 0 1 * *': 'Le 1er de chaque mois à minuit',
        '* * * * *': 'Toutes les minutes'
      }

      if (commonExpressions[cronExpression]) {
        return commonExpressions[cronExpression]
      }

      // Analyse des intervalles
      if (minute.startsWith('*/')) {
        const interval = parseInt(minute.substring(2))
        if (interval === 1) return 'Toutes les minutes'
        return `Toutes les ${interval} minutes`
      }

      if (hour.startsWith('*/')) {
        const interval = parseInt(hour.substring(2))
        return `Toutes les ${interval} heures`
      }

      // Analyse des heures spécifiques
      if (minute === '0' && hour !== '*' && day === '*' && month === '*' && dayOfWeek === '*') {
        const hourNum = parseInt(hour)
        if (hourNum === 0) return 'Tous les jours à minuit'
        if (hourNum < 12) return `Tous les jours à ${hourNum}h`
        if (hourNum === 12) return 'Tous les jours à midi'
        return `Tous les jours à ${hourNum}h`
      }

      // Analyse des jours de la semaine
      if (minute === '0' && hour === '0' && day === '*' && month === '*' && dayOfWeek !== '*') {
        const days = ['dimanche', 'lundi', 'mardi', 'mercredi', 'jeudi', 'vendredi', 'samedi']
        const dayIndex = parseInt(dayOfWeek)
        if (dayIndex >= 0 && dayIndex <= 6) {
          return `Tous les ${days[dayIndex]}s à minuit`
        }
      }

      // Analyse des jours du mois
      if (minute === '0' && hour === '0' && day !== '*' && month === '*' && dayOfWeek === '*') {
        const dayNum = parseInt(day)
        if (dayNum === 1) return 'Le 1er de chaque mois à minuit'
        if (dayNum === 15) return 'Le 15 de chaque mois à minuit'
        return `Le ${dayNum} de chaque mois à minuit`
      }

      // Pour les expressions plus complexes
      return 'Fréquence personnalisée'
    } catch {
      return 'Expression invalide'
    }
  }

  /**
   * Calcule la prochaine exécution d'une expression cron
   */
  static getNextExecution(cronExpression: string): string {
    try {
      const now = new Date()
      const parts = cronExpression.trim().split(/\s+/)
      if (parts.length !== 5) return 'Expression invalide'

      const [minute, hour, day, month, dayOfWeek] = parts

      // Toutes les minutes
      if (minute === '*' && hour === '*' && day === '*' && month === '*' && dayOfWeek === '*') {
        const nextTime = new Date(now)
        nextTime.setMinutes(now.getMinutes() + 1, 0, 0)
        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à ${nextTime.getHours().toString().padStart(2, '0')}:${nextTime.getMinutes().toString().padStart(2, '0')}`
      }

      // Toutes les X minutes
      if (minute.startsWith('*/') && hour === '*' && day === '*' && month === '*' && dayOfWeek === '*') {
        const interval = parseInt(minute.substring(2))
        const currentMinute = now.getMinutes()
        const nextMinute = Math.ceil(currentMinute / interval) * interval
        const nextTime = new Date(now)

        if (nextMinute >= 60) {
          nextTime.setHours(now.getHours() + 1, 0, 0, 0)
        } else {
          nextTime.setMinutes(nextMinute, 0, 0)
        }

        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à ${nextTime.getHours().toString().padStart(2, '0')}:${nextTime.getMinutes().toString().padStart(2, '0')}`
      }

      // Toutes les heures
      if (minute === '0' && hour === '*' && day === '*' && month === '*' && dayOfWeek === '*') {
        const nextTime = new Date(now)
        nextTime.setHours(now.getHours() + 1, 0, 0, 0)
        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à ${nextTime.getHours().toString().padStart(2, '0')}:00`
      }

      // Toutes les X heures
      if (minute === '0' && hour.startsWith('*/') && day === '*' && month === '*' && dayOfWeek === '*') {
        const interval = parseInt(hour.substring(2))
        const currentHour = now.getHours()
        const nextHour = Math.ceil(currentHour / interval) * interval
        const nextTime = new Date(now)

        if (nextHour >= 24) {
          nextTime.setDate(now.getDate() + 1)
          nextTime.setHours(0, 0, 0, 0)
        } else {
          nextTime.setHours(nextHour, 0, 0, 0)
        }

        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à ${nextTime.getHours().toString().padStart(2, '0')}:00`
      }

      // Heures spécifiques quotidiennes
      if (minute !== '*' && hour !== '*' && day === '*' && month === '*' && dayOfWeek === '*') {
        const scheduledHour = parseInt(hour)
        const scheduledMinute = parseInt(minute)
        const scheduledTime = new Date(now)
        scheduledTime.setHours(scheduledHour, scheduledMinute, 0, 0)

        if (scheduledTime <= now) {
          scheduledTime.setDate(scheduledTime.getDate() + 1)
        }

        return `Prochaine: ${scheduledTime.toLocaleDateString('fr-FR')} à ${scheduledHour.toString().padStart(2, '0')}:${scheduledMinute.toString().padStart(2, '0')}`
      }

      // Jours de la semaine spécifiques
      if (minute === '0' && hour === '0' && day === '*' && month === '*' && dayOfWeek !== '*') {
        const targetDay = parseInt(dayOfWeek)
        const currentDay = now.getDay()
        const daysUntil = (targetDay - currentDay + 7) % 7
        const nextTime = new Date(now)
        nextTime.setDate(now.getDate() + (daysUntil === 0 ? 7 : daysUntil))
        nextTime.setHours(0, 0, 0, 0)
        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à 00:00`
      }

      // Jours du mois spécifiques
      if (minute === '0' && hour === '0' && day !== '*' && month === '*' && dayOfWeek === '*') {
        const targetDay = parseInt(day)
        const nextTime = new Date(now.getFullYear(), now.getMonth(), targetDay, 0, 0, 0, 0)

        if (nextTime <= now) {
          nextTime.setMonth(now.getMonth() + 1)
          nextTime.setDate(targetDay)
        }

        return `Prochaine: ${nextTime.toLocaleDateString('fr-FR')} à 00:00`
      }

      return 'Calcul en cours...'
    } catch {
      return 'Expression invalide'
    }
  }

  /**
   * Valide une expression cron basique
   */
  static validateCronExpression(expression: string): boolean {
    if (!expression || expression.trim() === '') {
      return false
    }

    const parts = expression.trim().split(/\s+/)
    return parts.length === 5
  }

  /**
   * Génère des exemples d'expressions cron courantes
   */
  static getCommonExamples(): Array<{ expression: string; description: string; category: string }> {
    return [
      // Quotidiennes
      { expression: '0 0 * * *', description: 'Tous les jours à minuit', category: 'Quotidienne' },
      { expression: '0 6 * * *', description: 'Tous les jours à 6h', category: 'Quotidienne' },
      { expression: '0 12 * * *', description: 'Tous les jours à midi', category: 'Quotidienne' },
      { expression: '0 18 * * *', description: 'Tous les jours à 18h', category: 'Quotidienne' },

      // Horaires
      { expression: '0 * * * *', description: 'Toutes les heures', category: 'Horaire' },
      { expression: '0 */6 * * *', description: 'Toutes les 6 heures', category: 'Horaire' },
      { expression: '0 */12 * * *', description: 'Toutes les 12 heures', category: 'Horaire' },

      // Minutes
      { expression: '* * * * *', description: 'Toutes les minutes', category: 'Minute' },
      { expression: '*/5 * * * *', description: 'Toutes les 5 minutes', category: 'Minute' },
      { expression: '*/15 * * * *', description: 'Toutes les 15 minutes', category: 'Minute' },
      { expression: '*/30 * * * *', description: 'Toutes les 30 minutes', category: 'Minute' },

      // Hebdomadaires
      { expression: '0 0 * * 1', description: 'Tous les lundis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 2', description: 'Tous les mardis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 3', description: 'Tous les mercredis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 4', description: 'Tous les jeudis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 5', description: 'Tous les vendredis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 6', description: 'Tous les samedis à minuit', category: 'Hebdomadaire' },
      { expression: '0 0 * * 0', description: 'Tous les dimanches à minuit', category: 'Hebdomadaire' },

      // Mensuelles
      { expression: '0 0 1 * *', description: 'Le 1er de chaque mois à minuit', category: 'Mensuelle' },
      { expression: '0 0 15 * *', description: 'Le 15 de chaque mois à minuit', category: 'Mensuelle' }
    ]
  }
}