describe('Schedule Management', () => {
  const testEmail = `schedule.test.${Date.now()}@e2e.com`
  const testPassword = 'sT5!nK8@pQ3#vM9$'
  let testDatabaseId: number

  // Register user once before all tests
  before(() => {
    const baseUrl = Cypress.config('baseUrl') || 'http://localhost:3000'
    const apiBaseUrl = baseUrl.replace(':3000', ':8080').replace(':5173', ':8080')
    
    // Register the user
    cy.request({
      method: 'POST',
      url: `${apiBaseUrl}/auth/register`,
      body: {
        email: testEmail,
        password: testPassword,
        firstname: 'Schedule',
        lastname: 'Tester',
        confirm_password: testPassword
      },
      failOnStatusCode: false
    })
  })

  // Login before each test
  beforeEach(() => {
    const baseUrl = Cypress.config('baseUrl') || 'http://localhost:3000'
    const apiBaseUrl = baseUrl.replace(':3000', ':8080').replace(':5173', ':8080')
    
    // Login to get auth cookie
    cy.request({
      method: 'POST',
      url: `${apiBaseUrl}/auth/login`,
      body: {
        email: testEmail,
        password: testPassword
      }
    })
    
    // Verify authentication
    cy.getCookie('auth_token').should('exist')
  })

  describe('Schedule List View', () => {
    it('should display schedules page correctly', () => {
      cy.visit('/user/schedules')
      cy.url().should('include', '/user/schedules')
      
      // Check page title
      cy.contains('Mes sauvegardes planifiées').should('be.visible')
      
      // Check "Create" button exists
      cy.contains('button', 'Nouvelle planification').should('be.visible')
      
      // Check filter section
      cy.contains('label', 'Base de données:').should('be.visible')
    })

    it('should show empty state when no schedules', () => {
      cy.visit('/user/schedules')
      
      // Should show empty state message
      cy.contains('Aucune sauvegarde planifiée').should('be.visible')
      cy.contains('Créer votre première planification').should('be.visible')
    })
  })

  describe('Create Schedule', () => {
    before(() => {
      // Create a test database first
      const baseUrl = Cypress.config('baseUrl') || 'http://localhost:3000'
      const apiBaseUrl = baseUrl.replace(':3000', ':8080').replace(':5173', ':8080')
      
      cy.request({
        method: 'POST',
        url: `${apiBaseUrl}/auth/login`,
        body: {
          email: testEmail,
          password: testPassword
        }
      })
      
      cy.visit('/user/databases')
      cy.contains('button', 'Nouvelle base de données').click()
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Schedule Test DB')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('schedule_test_db')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('schedule_user')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('cX7!fR4@mK9#pL3$')
      cy.get('form').find('button[type="submit"]').click()
      cy.wait(2000)
    })

    it('should open create modal when clicking "Nouvelle planification"', () => {
      cy.visit('/user/schedules')
      
      cy.contains('button', 'Nouvelle planification').click()
      
      // Modal should be visible
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.contains('label', 'Nom').should('be.visible')
      cy.contains('label', 'Base de données').should('be.visible')
    })

    it('should create daily schedule successfully', () => {
      cy.visit('/user/schedules')
      
      cy.contains('button', 'Nouvelle planification').click()
      
      // Fill in the form
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Daily Backup Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1) // Select first database
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * *') // Daily at midnight
        
        // Submit the form
        cy.get('button[type="submit"]').click()
      })
      
      // Should show success
      cy.contains('Daily Backup Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should create weekly schedule successfully', () => {
      cy.visit('/user/schedules')
      
      cy.contains('button', 'Nouvelle planification').click()
      
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Weekly Backup Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1)
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * 1') // Weekly on Monday
        
        cy.get('button[type="submit"]').click()
      })
      
      cy.contains('Weekly Backup Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should cancel schedule creation', () => {
      cy.visit('/user/schedules')
      
      cy.contains('button', 'Nouvelle planification').click()
      
      // Fill some fields
      cy.contains('label', 'Nom').parent().find('input').type('Schedule to Cancel')
      
      // Click cancel
      cy.contains('button', 'Annuler').click()
      
      // Modal should close
      cy.contains('h2', 'Nouvelle planification').should('not.exist')
      
      // Schedule should not be created
      cy.contains('Schedule to Cancel').should('not.exist')
    })
  })

  describe('Schedule Details', () => {
    beforeEach(() => {
      // Create a schedule for testing
      cy.visit('/user/schedules')
      cy.contains('button', 'Nouvelle planification').click()
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Details Test Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1)
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * *') // Daily
        cy.get('button[type="submit"]').click()
      })
      cy.contains('Details Test Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should display schedule information correctly', () => {
      cy.visit('/user/schedules')
      
      // Check that schedule card shows correct information
      cy.contains('Details Test Schedule').should('be.visible')
      
      // Check card content
      cy.contains('Details Test Schedule').parents('.bg-white').within(() => {
        cy.contains('Base de données').should('be.visible')
        cy.contains('Fréquence:').should('be.visible')
        cy.contains('Prochaine exécution:').should('be.visible')
        cy.contains('Créé le:').should('be.visible')
      })
    })

    it('should show active/inactive status badge', () => {
      cy.visit('/user/schedules')
      
      cy.contains('Details Test Schedule').parents('.bg-white').within(() => {
        // Should have either "Actif" or "Inactif" badge
        cy.get('span').contains(/Actif|Inactif/).should('be.visible')
      })
    })

    it('should show action buttons for schedule', () => {
      cy.visit('/user/schedules')
      
      cy.contains('Details Test Schedule').parents('.bg-white').within(() => {
        // Edit button (pencil icon)
        cy.get('button').eq(0).should('exist')
        
        // Delete button (trash icon)
        cy.get('button').eq(1).should('exist')
        
        // Toggle button (Activer/Désactiver)
        cy.contains('button', /Activer|Désactiver/).should('be.visible')
      })
    })
  })

  describe('Toggle Schedule', () => {
    beforeEach(() => {
      cy.visit('/user/schedules')
      cy.contains('button', 'Nouvelle planification').click()
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Toggle Test Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1)
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * *') // Daily
        cy.get('button[type="submit"]').click()
      })
      cy.contains('Toggle Test Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should toggle schedule activation status', () => {
      cy.visit('/user/schedules')
      
      cy.contains('Toggle Test Schedule').parents('.bg-white').within(() => {
        // Get initial status
        cy.get('span').contains(/Actif|Inactif/).invoke('text').then((initialStatus) => {
          // Click toggle button
          cy.contains('button', /Activer|Désactiver/).click()
          
          cy.wait(1000)
          
          // Status should have changed
          cy.get('span').contains(/Actif|Inactif/).should('not.contain', initialStatus)
        })
      })
    })
  })

  describe('Update Schedule', () => {
    beforeEach(() => {
      cy.visit('/user/schedules')
      cy.contains('button', 'Nouvelle planification').click()
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Update Test Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1)
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * *') // Daily
        cy.get('button[type="submit"]').click()
      })
      cy.contains('Update Test Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should open edit modal when clicking edit button', () => {
      cy.visit('/user/schedules')
      
      // Click edit button (first button in the card, pencil icon)
      cy.contains('Update Test Schedule').parents('.bg-white').find('button').eq(0).click()
      
      // Edit modal should appear
      cy.contains('h2', 'Modifier la planification').should('be.visible')
    })

    it('should update schedule name', () => {
      cy.visit('/user/schedules')
      
      // Click edit button
      cy.contains('Update Test Schedule').parents('.bg-white').find('button').eq(0).click()
      
      // Change the name
      cy.contains('label', 'Nom').parent().find('input').clear().type('Updated Schedule Name')
      
      // Submit
      cy.contains('button', 'Modifier').click()
      
      cy.wait(2000)
      
      // Should show updated name
      cy.contains('Updated Schedule Name', { timeout: 10000 }).should('be.visible')
      
      // Old name should not exist in the specific card
      cy.contains('Updated Schedule Name').parents('.bg-white').within(() => {
        cy.get('h3').should('not.contain', 'Update Test Schedule')
      })
    })

    it('should cancel edit without saving', () => {
      cy.visit('/user/schedules')
      
      // Click edit button
      cy.contains('Update Test Schedule').parents('.bg-white').find('button').eq(0).click()
      
      // Change the name
      cy.contains('label', 'Nom').parent().find('input').clear().type('Cancelled Name')
      
      // Cancel
      cy.contains('button', 'Annuler').click()
      
      // Name should not change
      cy.contains('Update Test Schedule').should('be.visible')
      cy.contains('Cancelled Name').should('not.exist')
    })
  })

  describe('Delete Schedule', () => {
    beforeEach(() => {
      cy.visit('/user/schedules')
      cy.contains('button', 'Nouvelle planification').click()
      cy.contains('h2', 'Nouvelle planification').should('be.visible')
      cy.get('form').within(() => {
        cy.contains('label', 'Nom').parent().find('input').type('Delete Test Schedule')
        cy.contains('label', 'Base de données').parent().find('select').select(1)
        cy.contains('label', 'Fréquence').parent().find('select').select('0 0 * * *') // Daily
        cy.get('button[type="submit"]').click()
      })
      cy.contains('Delete Test Schedule', { timeout: 10000 }).should('be.visible')
    })

    it('should delete schedule with confirmation', () => {
      cy.visit('/user/schedules')
      
      // Stub the confirm dialog to automatically accept
      cy.window().then((win) => {
        cy.stub(win, 'confirm').returns(true)
      })
      
      // Click delete button (second button in the card, trash icon)
      cy.contains('Delete Test Schedule').parents('.bg-white').find('button').eq(1).click()
      
      // Schedule should be removed
      cy.contains('Delete Test Schedule', { timeout: 10000 }).should('not.exist')
    })

    it('should cancel deletion', () => {
      cy.visit('/user/schedules')
      
      // Stub the confirm dialog to automatically reject
      cy.window().then((win) => {
        cy.stub(win, 'confirm').returns(false)
      })
      
      // Click delete button
      cy.contains('Delete Test Schedule').parents('.bg-white').find('button').eq(1).click()
      
      // Schedule should still exist
      cy.contains('Delete Test Schedule').should('be.visible')
    })
  })

  describe('Filter Schedules', () => {
    it('should filter schedules by database', () => {
      cy.visit('/user/schedules')
      
      // Check that database filter exists
      cy.contains('label', 'Base de données:').should('be.visible')
      cy.get('select').should('be.visible')
      
      // Select "Toutes les bases"
      cy.get('select').select('Toutes les bases')
    })
  })
})
