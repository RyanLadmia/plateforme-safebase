// ***********************************************
// Custom Cypress Commands
// ***********************************************

/**
 * Simulate authenticated user by setting auth cookie directly
 * This bypasses the need for registration and login
 * @example cy.authenticateUser()
 */
Cypress.Commands.add('authenticateUser', (email?: string, password?: string) => {
  // Use default test user if not provided
  const userEmail = email || 'cypress.test@e2e.com'
  const userPassword = password || 'xK9#mQ2$vL7@wP4!nR8%'
  
  // Get the base URL without /api for auth routes
  const baseUrl = Cypress.config('baseUrl') || 'http://localhost:3000'
  const apiBaseUrl = baseUrl.replace(':3000', ':8080') // Frontend port 3000, backend port 8080
  
  cy.session([userEmail, userPassword], () => {
    // Try to login first
    cy.request({
      method: 'POST',
      url: `${apiBaseUrl}/auth/login`,
      body: {
        email: userEmail,
        password: userPassword
      },
      failOnStatusCode: false
    }).then((response) => {
      if (response.status === 200) {
        // Login successful, cookie is set automatically
        return
      }
      
      // If login fails (user doesn't exist), register first
      if (response.status === 401 || response.status === 400) {
        cy.request({
          method: 'POST',
          url: `${apiBaseUrl}/auth/register`,
          body: {
            email: userEmail,
            password: userPassword,
            firstname: 'Cypress',
            lastname: 'Test',
            confirm_password: userPassword
          }
        }).then(() => {
          // Now login
          cy.request({
            method: 'POST',
            url: `${apiBaseUrl}/auth/login`,
            body: {
              email: userEmail,
              password: userPassword
            }
          })
        })
      }
    })
    
    // Verify authentication worked
    cy.getCookie('auth_token').should('exist')
  }, {
    cacheAcrossSpecs: true
  })
})

/**
 * Login command - Authenticates a user via UI
 * @deprecated Use cy.authenticateUser() instead for faster tests
 */
Cypress.Commands.add('login', (email: string, password: string) => {
  cy.session([email, password], () => {
    cy.visit('/login')
    
    // S'assurer qu'on est sur l'onglet connexion
    cy.contains('button', 'Connexion').click()
    
    cy.get('input#login-email').type(email)
    cy.get('input#login-password').type(password)
    cy.get('button[type="submit"]').click()
    
    // Wait for successful login (redirige vers dashboard)
    cy.url({ timeout: 10000 }).should('match', /dashboard/)
    
    // Vérifier qu'un cookie de session existe
    cy.getCookies().should('have.length.at.least', 1)
  })
})

/**
 * Logout command - Logs out the current user
 */
Cypress.Commands.add('logout', () => {
  // Chercher le bouton de déconnexion (adapter selon votre UI)
  cy.contains(/déconnexion|logout|se déconnecter/i, { timeout: 5000 }).click()
  cy.url({ timeout: 5000 }).should('include', '/login')
  
  // Vérifier que les cookies sont supprimés
  cy.getCookie('token').should('be.null')
})

/**
 * Register user command - Creates a new user account
 * Note: L'inscription NE redirige PAS automatiquement vers le dashboard
 * @deprecated Use cy.authenticateUser() instead for faster tests
 */
Cypress.Commands.add('registerUser', (userData) => {
  cy.visit('/login')
  
  // Cliquer sur l'onglet Inscription
  cy.contains('button', 'Inscription').click()
  
  cy.get('input#register-firstname').type(userData.firstname)
  cy.get('input#register-lastname').type(userData.lastname)
  cy.get('input#register-email').type(userData.email)
  cy.get('input#register-password').type(userData.password)
  cy.get('input#register-confirm-password').type(userData.password)
  
  cy.get('button[type="submit"]').click()
  
  // Attendre le message de succès (PAS de redirection automatique)
  cy.contains(/inscription réussie|compte créé|utilisateur.*créé/i, { timeout: 10000 })
    .should('be.visible')
})

/**
 * NOTE: Les commandes createDatabase, createSchedule, et deleteAllTestData
 * ont été supprimées car l'application utilise des cookies HTTP-only
 * au lieu de localStorage pour l'authentification.
 * 
 * Pour créer des ressources dans les tests, utilisez l'UI directement :
 * 
 * Exemple pour créer une base de données :
 *   cy.visit('/user/databases')
 *   cy.contains('button', 'Nouvelle base de données').click()
 *   // ... remplir le formulaire
 *   cy.contains('button', 'Créer').click()
 */

/**
 * Check accessibility - Basic accessibility checks
 */
Cypress.Commands.add('checkAccessibility', () => {
  // Check for basic accessibility features
  cy.get('main').should('exist')
  cy.get('nav').should('exist')
  
  // Check for semantic HTML
  cy.get('button').each(($btn) => {
    cy.wrap($btn).should('have.attr', 'type')
  })
  
  // Check for alt text on images
  cy.get('img').each(($img) => {
    cy.wrap($img).should('have.attr', 'alt')
  })
})

/**
 * Clean up test users - Removes all users with @e2e.com email domain
 * This should be called after test suites to clean up the database
 */
Cypress.Commands.add('cleanupTestUsers', () => {
  const apiBaseUrl = Cypress.config('baseUrl')?.replace(':3000', ':8080') || 'http://localhost:8080'
  
  cy.request({
    method: 'POST',
    url: `${apiBaseUrl}/api/test/cleanup-users`,
    failOnStatusCode: false
  }).then((response) => {
    if (response.status === 200) {
      cy.log(`Cleaned up ${response.body.deleted_count || 0} test users`)
    } else {
      cy.log('Cleanup endpoint not available or failed')
    }
  })
})

export {}

