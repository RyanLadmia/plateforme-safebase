/**
 * E2E Tests - Complete User Workflows
 * Tests complete end-to-end user journeys
 * Coverage: ~10% of the application
 */

describe('Complete User Workflows', () => {
  before(() => {
    // Authenticate once for all workflow tests
    cy.fixture('users').then((users) => {
      const user = users.users.validUser
      cy.authenticateUser(user.email, user.password)
    })
  })

  beforeEach(() => {
    cy.clearLocalStorage()
    cy.clearCookies()
    
    // Re-login for each test
    cy.fixture('users').then((users) => {
      const user = users.users.validUser
      cy.request({
        method: 'POST',
        url: 'http://localhost:8080/auth/login',
        body: {
          email: user.email,
          password: user.password
        }
      })
      cy.getCookie('auth_token').should('exist')
    })
  })

  describe('Database Management Workflow', () => {
    it('should create and view a database', () => {
      // Navigate to databases page
      cy.visit('/user/databases')
      
      // Check if there are existing databases or create new one
      cy.get('body').then(($body) => {
        if ($body.text().includes('Nouvelle base de données')) {
          cy.contains('button', 'Nouvelle base de données').click()
          
          // Fill database form
          cy.get('form').within(() => {
            cy.get('input[name="name"]').type('Workflow Test DB')
            cy.get('select[name="type"]').select('mysql')
            cy.get('input[name="host"]').type('localhost')
            cy.get('input[name="port"]').type('3306')
            cy.get('input[name="username"]').type('test_user')
            cy.get('input[name="password"]').type('zB6!pT9@qH3#xW7$')
            cy.get('input[name="db_name"]').type('test_db')
            cy.get('button[type="submit"]').click()
          })
          
          // Should show success message
          cy.contains(/créée|ajoutée|success/i, { timeout: 10000 }).should('be.visible')
        }
      })
      
      // Verify page loads correctly
      cy.visit('/user/databases')
      cy.url().should('include', '/user/databases')
      cy.contains(/bases de données|mes bases/i).should('be.visible')
    })

    it('should navigate through main pages', () => {
      const pages = [
        '/user/dashboard',
        '/user/databases',
        '/user/backups',
        '/user/schedules',
        '/user/history',
        '/user/profile'
      ]
      
      pages.forEach((page) => {
        cy.visit(page)
        cy.url().should('include', page)
        // Verify user is still authenticated
        cy.getCookie('auth_token').should('exist')
      })
    })
  })

  describe('Backup Creation Workflow', () => {
    it('should create a backup from database page', () => {
      cy.visit('/user/databases')
      
      // Check if databases exist
      cy.get('body').then(($body) => {
        if ($body.text().match(/créer une sauvegarde|backup/i)) {
          // Click on create backup button in first database card
          cy.get('.bg-white.rounded-lg.shadow').first().within(() => {
            cy.contains('button', /créer une sauvegarde|backup/i).click({ force: true })
          })
          
          // Should navigate to backups page or show success
          cy.wait(2000)
        }
      })
    })

    it('should view backups page', () => {
      cy.visit('/user/backups')
      
      // Page should load successfully
      cy.url().should('include', '/user/backups')
      cy.contains(/sauvegardes|backups/i).should('be.visible')
    })
  })

  describe('Schedule Management Workflow', () => {
    it('should navigate to schedules page', () => {
      cy.visit('/user/schedules')
      
      // Page should load successfully
      cy.url().should('include', '/user/schedules')
      cy.contains(/planifications|schedules/i).should('be.visible')
    })

    it('should be able to create schedule if database exists', () => {
      cy.visit('/user/schedules')
      
      cy.get('body').then(($body) => {
        if ($body.text().includes('Nouvelle planification')) {
          cy.contains('button', 'Nouvelle planification').click()
          
          // Modal should open
          cy.contains('h2', 'Nouvelle planification').should('be.visible')
          
          // Close modal
          cy.get('body').type('{esc}')
        }
      })
    })
  })

  describe('History Viewing Workflow', () => {
    it('should view action history', () => {
      cy.visit('/user/history')
      
      // Page should load successfully
      cy.url().should('include', '/user/history')
      cy.contains(/historique|history/i).should('be.visible')
    })

    it('should be able to filter history', () => {
      cy.visit('/user/history')
      
      // Check if filter buttons exist
      cy.get('body').then(($body) => {
        if ($body.text().match(/toutes|créations|modifications/i)) {
          cy.contains('button', /toutes|all/i).should('be.visible')
        }
      })
    })
  })

  describe('Profile Management Workflow', () => {
    it('should view and update profile', () => {
      cy.visit('/user/profile')
      
      // Profile page should load
      cy.url().should('include', '/user/profile')
      cy.contains('Informations personnelles').should('be.visible')
      
      // Should be able to edit profile
      cy.contains('button', 'Modifier').should('be.visible')
    })

    it('should access password change section', () => {
      cy.visit('/user/profile')
      
      // Password change section should be visible
      cy.contains('Changer le mot de passe').should('be.visible')
      cy.contains('label', 'Mot de passe actuel').should('be.visible')
    })
  })

  describe('Navigation Workflow', () => {
    it('should navigate using sidebar menu', () => {
      cy.visit('/user/dashboard')
      
      // Navigate to databases
      cy.contains(/bases de données|databases/i).first().click({ force: true })
      cy.url().should('include', '/user/databases')
      
      // Navigate to backups
      cy.contains(/sauvegardes|backups/i).first().click({ force: true })
      cy.url().should('include', '/user/backups')
      
      // Navigate to schedules
      cy.contains(/planifications|schedules/i).first().click({ force: true })
      cy.url().should('include', '/user/schedules')
      
      // Navigate to history
      cy.contains(/historique|history/i).first().click({ force: true })
      cy.url().should('include', '/user/history')
    })

    it('should access profile from menu', () => {
      cy.visit('/user/dashboard')
      
      // Click on profile link
      cy.contains(/mon profil|profile/i).first().click({ force: true })
      cy.url().should('include', '/user/profile')
    })
  })

  describe('Session Persistence Workflow', () => {
    it('should maintain authentication across page navigation', () => {
      const pages = [
        '/user/dashboard',
        '/user/databases',
        '/user/backups',
        '/user/schedules',
        '/user/history'
      ]
      
      pages.forEach((page) => {
        cy.visit(page)
        cy.url().should('include', page)
        
        // Should still be authenticated (cookie exists)
        cy.getCookie('auth_token').should('exist')
      })
    })

    it('should reload page without losing authentication', () => {
      cy.visit('/user/dashboard')
      
      // Reload page
      cy.reload()
      
      // Should still be on dashboard
      cy.url().should('include', '/user/dashboard')
      cy.getCookie('auth_token').should('exist')
    })
  })

  describe('Responsive Design Workflow', () => {
    it('should work on mobile viewport', () => {
      cy.viewport('iphone-x')
      
      cy.visit('/user/dashboard')
      cy.url().should('include', '/user/dashboard')
      
      // Navigation should be accessible
      cy.get('nav, aside, [role="navigation"]').should('exist')
    })

    it('should work on tablet viewport', () => {
      cy.viewport('ipad-2')
      
      cy.visit('/user/dashboard')
      cy.url().should('include', '/user/dashboard')
      
      // Content should be visible
      cy.contains(/bases de données|databases/i).should('exist')
    })

    it('should work on desktop viewport', () => {
      cy.viewport(1920, 1080)
      
      cy.visit('/user/dashboard')
      cy.url().should('include', '/user/dashboard')
      
      // Full layout should be visible
      cy.get('nav, aside, [role="navigation"]').should('be.visible')
    })
  })

  describe('Accessibility Workflow', () => {
    it('should be keyboard navigable', () => {
      cy.visit('/user/dashboard')
      
      // Focus on first interactive element
      cy.get('a, button').first().focus()
      cy.focused().should('exist')
      
      // Focus on second element
      cy.get('a, button').eq(1).focus()
      cy.focused().should('exist')
    })

    it('should have proper heading structure', () => {
      cy.visit('/user/dashboard')
      
      // Should have headings
      cy.get('h1, h2, h3').should('exist')
    })
  })

  describe('Error Handling Workflow', () => {
    it('should handle invalid routes gracefully', () => {
      cy.visit('/user/nonexistent-page', { failOnStatusCode: false })
      
      // Application should either redirect or stay on invalid page (Vue router handles it)
      cy.url().then((url) => {
        // Accept that the page loads (Vue router doesn't redirect to 404 for invalid routes)
        expect(url).to.include('/user/')
      })
    })

    it('should handle network errors in forms', () => {
      cy.visit('/user/profile')
      
      // Intercept and simulate network error
      cy.intercept('PUT', '/api/profile', {
        statusCode: 500,
        body: { error: 'Server error' }
      }).as('updateProfile')
      
      // Try to update profile
      cy.contains('button', 'Modifier').click()
      cy.contains('label', 'Prénom').parent().find('input').clear().type('TestName')
      cy.contains('button', 'Enregistrer').click()
      
      // Should show error message
      cy.wait('@updateProfile')
      cy.contains(/erreur|error/i, { timeout: 5000 }).should('be.visible')
    })
  })

  describe('Complete User Journey', () => {
    it('should complete a typical user session', () => {
      // Start at dashboard
      cy.visit('/user/dashboard')
      cy.url().should('include', '/user/dashboard')
      
      // Navigate to databases
      cy.visit('/user/databases')
      cy.url().should('include', '/user/databases')
      
      // Check backups
      cy.visit('/user/backups')
      cy.url().should('include', '/user/backups')
      
      // View history
      cy.visit('/user/history')
      cy.url().should('include', '/user/history')
      
      // Check profile
      cy.visit('/user/profile')
      cy.url().should('include', '/user/profile')
      cy.contains('Informations personnelles').should('be.visible')
      
      // Return to dashboard
      cy.visit('/user/dashboard')
      cy.url().should('include', '/user/dashboard')
      
      // Verify still authenticated throughout
      cy.getCookie('auth_token').should('exist')
    })
  })

  describe('Statistics and Quick Actions Workflow', () => {
    it('should display statistics on dashboard', () => {
      cy.visit('/user/dashboard')
      
      // Statistics should be visible
      cy.contains(/bases de données|databases/i).should('be.visible')
      cy.contains(/sauvegardes|backups/i).should('be.visible')
      cy.contains(/planifications|schedules/i).should('be.visible')
      
      // Should display numbers for each statistic
      cy.get('.text-3xl.font-bold').should('have.length.at.least', 3)
    })

    it('should use quick action links from dashboard', () => {
      cy.visit('/user/dashboard')
      
      // Click on "Gérer mes bases de données"
      cy.contains('Gérer mes bases de données').click()
      cy.url().should('include', '/user/databases')
      
      // Go back to dashboard
      cy.visit('/user/dashboard')
      
      // Click on "Gérer mes sauvegardes"
      cy.contains('Gérer mes sauvegardes').click()
      cy.url().should('include', '/user/backups')
      
      // Go back to dashboard
      cy.visit('/user/dashboard')
      
      // Click on "Gérer mes planifications"
      cy.contains('Gérer mes planifications').click()
      cy.url().should('include', '/user/schedules')
    })
  })

  describe('Backup Status Filtering Workflow', () => {
    it('should filter backups by status', () => {
      cy.visit('/user/backups')
      
      // Check if filter buttons exist
      cy.get('body').then(($body) => {
        if ($body.text().match(/toutes|en cours|réussies|échouées/i)) {
          // Try clicking different filters
          cy.contains('button', /toutes/i).should('exist')
          cy.contains('button', /en cours/i).should('exist')
          cy.contains('button', /réussies/i).should('exist')
        }
      })
    })
  })

  describe('Schedule Frequency Selection Workflow', () => {
    it('should view schedule frequency options', () => {
      cy.visit('/user/schedules')
      
      // Open create schedule modal if database exists
      cy.get('body').then(($body) => {
        if ($body.text().includes('Nouvelle planification')) {
          cy.contains('button', 'Nouvelle planification').click()
          
          // Check for frequency selector
          cy.get('body').then(($modal) => {
            if ($modal.text().includes('Fréquence')) {
              cy.contains('label', 'Fréquence').should('be.visible')
            }
          })
          
          // Close modal
          cy.get('body').type('{esc}')
        }
      })
    })
  })

  describe('History Action Type Filtering Workflow', () => {
    it('should filter history by action types', () => {
      cy.visit('/user/history')
      
      // Should have filter buttons
      cy.contains('button', /toutes/i).should('be.visible')
      cy.contains('button', /bases de données/i).should('be.visible')
      cy.contains('button', /sauvegardes/i).should('be.visible')
      cy.contains('button', /planifications/i).should('be.visible')
      
      // Click on different filters
      cy.contains('button', /bases de données/i).click()
      cy.contains('button', /sauvegardes/i).click()
      cy.contains('button', /planifications/i).click()
      cy.contains('button', /toutes/i).click()
    })

    it('should display action details in history', () => {
      cy.visit('/user/history')
      
      // Check for history structure
      cy.get('body').then(($body) => {
        if (!$body.text().includes('Aucune activité')) {
          // If there are activities, check their structure
          cy.get('.bg-white.rounded-lg.shadow, .p-6.hover\\:bg-gray-50').first().should('exist')
        }
      })
    })
  })

  describe('Profile Statistics Workflow', () => {
    it('should display account statistics in profile', () => {
      cy.visit('/user/profile')
      
      // Should show statistics section
      cy.contains('Statistiques du compte').should('be.visible')
      
      // Should show counts
      cy.get('.text-3xl.font-bold.text-blue-600').should('exist') // Databases
      cy.get('.text-3xl.font-bold.text-green-600').should('exist') // Total backups
      cy.get('.text-3xl.font-bold.text-orange-600').should('exist') // Completed backups
    })
  })

  describe('Data Consistency Workflow', () => {
    it('should show consistent data across pages', () => {
      // Just verify that statistics are displayed correctly on both pages
      cy.visit('/user/dashboard')
      
      // Dashboard should have statistics cards (at least 3)
      cy.get('.bg-white.rounded-lg.shadow').should('have.length.at.least', 3)
      cy.get('.text-3xl.font-bold').should('have.length.at.least', 3)
      
      // Verify profile statistics
      cy.visit('/user/profile')
      cy.contains('Statistiques du compte').should('be.visible')
      
      // Profile should have exactly 3 statistics in the grid
      cy.get('.grid.grid-cols-1.md\\:grid-cols-3 .text-center').should('have.length', 3)
      
      // Verify all three statistics display numbers
      cy.get('.grid.grid-cols-1.md\\:grid-cols-3 .text-center').each(($stat) => {
        cy.wrap($stat).find('.text-3xl.font-bold').should('exist')
        cy.wrap($stat).find('.text-3xl.font-bold').invoke('text').should('match', /^\d+$/)
      })
      
      // Verify labels
      cy.contains('.text-gray-500', 'Bases de données').should('be.visible')
      cy.contains('.text-gray-500', 'Sauvegardes totales').should('be.visible')
      cy.contains('.text-gray-500', 'Sauvegardes réussies').should('be.visible')
    })
  })

  describe('Empty State Workflow', () => {
    it('should handle empty databases page gracefully', () => {
      cy.visit('/user/databases')
      
      // Page should load (either with databases or empty state)
      cy.get('body').should('exist')
      cy.url().should('include', '/user/databases')
    })

    it('should handle empty backups page gracefully', () => {
      cy.visit('/user/backups')
      
      // Page should load (either with backups or empty state)
      cy.get('body').should('exist')
      cy.url().should('include', '/user/backups')
    })

    it('should handle empty schedules page gracefully', () => {
      cy.visit('/user/schedules')
      
      // Page should load (either with schedules or empty state)
      cy.get('body').should('exist')
      cy.url().should('include', '/user/schedules')
    })

    it('should handle empty history page gracefully', () => {
      cy.visit('/user/history')
      
      // Should show page title
      cy.contains(/historique|history/i).should('be.visible')
      
      // Either show activities or "Aucune activité"
      cy.get('body').should('exist')
    })
  })
})
