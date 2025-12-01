// Configuration d'Axios pour l'application
import axios from 'axios'

const API_BASE_URL: string = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Instance Axios configurée
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // CRUCIAL: Envoie automatiquement les cookies HTTP-only avec chaque requête
  headers: {
    'Content-Type': 'application/json',
  },
})

// Note de sécurité : 
// Le token JWT est stocké dans un cookie HTTP-only par le backend.
// Il est automatiquement envoyé avec chaque requête grâce à withCredentials: true.
// Cela protège contre les attaques XSS car JavaScript ne peut pas accéder au cookie.

// Intercepteur de réponse pour gérer les erreurs globalement
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Extraire le message d'erreur du backend
    const message = error.response?.data?.error || error.message || 'Une erreur est survenue'
    
    // Créer une erreur avec le message approprié
    const customError = new Error(message)
    
    // Ajouter des informations supplémentaires
    ;(customError as any).status = error.response?.status
    ;(customError as any).data = error.response?.data
    
    return Promise.reject(customError)
  }
)

export default apiClient
