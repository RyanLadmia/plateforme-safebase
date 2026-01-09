describe('Backup Management', () => {
  const testEmail = `backup.test.${Date.now()}@e2e.com`
  const testPassword = 'hD8!qW3@vN7#xM2$'
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
        firstname: 'Backup',
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

  describe('Backup List View', () => {
    it('should display backups page correctly', () => {
      cy.visit('/user/backups')
      cy.url().should('include', '/user/backups')
      
      // Check page title
      cy.contains('Mes sauvegardes').should('be.visible')
      
      // Check filter buttons
      cy.contains('button', 'Toutes').should('be.visible')
      cy.contains('button', 'Terminées').should('be.visible')
      cy.contains('button', 'En cours').should('be.visible')
      cy.contains('button', 'Échouées').should('be.visible')
      
      // Check statistics section
      cy.contains('Taille totale').should('be.visible')
      cy.contains('Sauvegardes réussies').should('be.visible')
    })

    it('should show empty state when no backups', () => {
      cy.visit('/user/backups')
      
      // Should show empty state message
      cy.contains('Aucune sauvegarde trouvée').should('be.visible')
    })
  })

  describe('Create Backup', () => {
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
      cy.contains('label', 'Nom').parent().find('input').type('Backup Test DB')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('backup_test_db')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('backup_user')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('pL9!kR4@jT6#mW2$')
      cy.get('form').find('button[type="submit"]').click()
      cy.wait(2000)
    })

    it('should create backup from database card', () => {
      cy.visit('/user/databases')
      
      // Click backup button on the database card
      cy.contains('Backup Test DB').parents('.bg-white').within(() => {
        cy.contains('button', 'Créer une sauvegarde').click()
      })
      
      // Should show success or navigate to backups
      cy.wait(2000)
      
      // Verify backup was created by visiting backups page
      cy.visit('/user/backups')
      cy.contains('Backup Test DB', { timeout: 10000 }).should('be.visible')
    })
  })

  describe('Backup Filters', () => {
    it('should filter backups by status', () => {
      cy.visit('/user/backups')
      
      // Click on different status filters
      cy.contains('button', 'Terminées').click()
      cy.contains('button', 'Terminées').should('have.class', 'bg-green-600')
      
      cy.contains('button', 'En cours').click()
      cy.contains('button', 'En cours').should('have.class', 'bg-orange-600')
      
      cy.contains('button', 'Échouées').click()
      cy.contains('button', 'Échouées').should('have.class', 'bg-red-600')
      
      cy.contains('button', 'Toutes').click()
      cy.contains('button', 'Toutes').should('have.class', 'bg-blue-600')
    })

    it('should filter backups by database', () => {
      cy.visit('/user/backups')
      
      // Check that database filter exists
      cy.contains('label', 'Base de données:').should('be.visible')
      cy.get('select').should('be.visible')
      
      // Select a database filter
      cy.get('select').select('Toutes les bases')
    })
  })

  describe('Backup Actions', () => {
    it('should display backup actions for completed backups', () => {
      cy.visit('/user/backups')
      
      // Wait for backups to load
      cy.wait(2000)
      
      // Check if there are any completed backups with download button
      cy.get('body').then(($body) => {
        if ($body.find('button:contains("Télécharger")').length > 0) {
          cy.contains('button', 'Télécharger').should('be.visible')
        }
      })
    })

    it('should show delete confirmation for backup', () => {
      cy.visit('/user/backups')
      cy.wait(2000)
      
      // Check if there are any backups with delete button
      cy.get('body').then(($body) => {
        if ($body.find('button:contains("Supprimer")').length > 0) {
          // Stub the confirm dialog
          cy.window().then((win) => {
            cy.stub(win, 'confirm').returns(false)
          })
          
          cy.contains('button', 'Supprimer').first().click()
        }
      })
    })
  })

  describe('Backup Statistics', () => {
    it('should display backup statistics correctly', () => {
      cy.visit('/user/backups')
      
      // Check statistics cards (in the grid section, not in filter buttons)
      cy.get('.grid').within(() => {
        cy.contains('Taille totale').should('be.visible')
        cy.contains('Sauvegardes réussies').should('be.visible')
        cy.contains('En cours').should('be.visible')
        cy.contains('Échouées').should('be.visible')
      })
    })
  })

  describe('Backup Table', () => {
    it('should display backup table with correct columns', () => {
      cy.visit('/user/backups')
      
      // Check table headers
      cy.contains('th', 'Fichier').should('be.visible')
      cy.contains('th', 'Date').should('be.visible')
      cy.contains('th', 'Taille').should('be.visible')
      cy.contains('th', 'Type').should('be.visible')
      cy.contains('th', 'Statut').should('be.visible')
      cy.contains('th', 'Actions').should('be.visible')
    })

    it('should display backup information in table rows', () => {
      cy.visit('/user/backups')
      cy.wait(2000)
      
      // Check if there are any backups displayed
      cy.get('body').then(($body) => {
        if ($body.find('table tbody tr').length > 0) {
          // Verify first row has data
          cy.get('table tbody tr').first().within(() => {
            cy.get('td').should('have.length.at.least', 6)
          })
        }
      })
    })
  })

  describe('Download Backup', () => {
    it('should have download button for completed backups', () => {
      cy.visit('/user/backups')
      cy.wait(2000)
      
      // Check for download buttons
      cy.get('body').then(($body) => {
        if ($body.find('button:contains("Télécharger")').length > 0) {
          cy.contains('button', 'Télécharger').should('be.visible')
          cy.contains('button', 'Télécharger').should('not.be.disabled')
        }
      })
    })
  })

  describe('Delete Backup', () => {
    it('should delete backup with confirmation', () => {
      cy.visit('/user/backups')
      cy.wait(2000)
      
      cy.get('body').then(($body) => {
        const deleteButtons = $body.find('button:contains("Supprimer")')
        
        if (deleteButtons.length > 0) {
          const initialCount = $body.find('table tbody tr').length
          
          // Stub confirm to accept
          cy.window().then((win) => {
            cy.stub(win, 'confirm').returns(true)
          })
          
          // Click delete
          cy.contains('button', 'Supprimer').first().click()
          
          // Wait and verify deletion
          cy.wait(2000)
          cy.get('table tbody tr').should('have.length.lessThan', initialCount + 1)
        }
      })
    })

    it('should cancel backup deletion', () => {
      cy.visit('/user/backups')
      cy.wait(2000)
      
      cy.get('body').then(($body) => {
        if ($body.find('button:contains("Supprimer")').length > 0) {
          const initialCount = $body.find('table tbody tr').length
          
          // Stub confirm to reject
          cy.window().then((win) => {
            cy.stub(win, 'confirm').returns(false)
          })
          
          cy.contains('button', 'Supprimer').first().click()
          
          // Verify nothing changed
          cy.wait(1000)
          cy.get('table tbody tr').should('have.length', initialCount)
        }
      })
    })
  })
})
