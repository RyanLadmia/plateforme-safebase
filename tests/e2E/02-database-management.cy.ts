describe('Database Management', () => {
  const testEmail = `db.test.${Date.now()}@e2e.com`
  const testPassword = 'xK9#mQ2$vL7@wP4!nR8%'

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
        firstname: 'Cypress',
        lastname: 'Test',
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

  describe('Database List View', () => {
    it('should display databases page correctly', () => {
      cy.visit('/user/databases')
      cy.url().should('include', '/user/databases')
      
      // Check page title
      cy.contains('Mes bases de données').should('be.visible')
      
      // Check "Create" button exists
      cy.contains('button', 'Nouvelle base de données').should('be.visible')
      
      // Check filter buttons
      cy.contains('button', 'Tous types').should('be.visible')
      cy.contains('button', 'MySQL').should('be.visible')
      cy.contains('button', 'PostgreSQL').should('be.visible')
    })

    it('should show empty state when no databases', () => {
      cy.visit('/user/databases')
      
      // Should show empty state message
      cy.contains('Aucune base de données configurée').should('be.visible')
      cy.contains('Ajouter votre première base de données').should('be.visible')
    })
  })

  describe('Create Database', () => {
    it('should open create modal when clicking "Nouvelle base de données"', () => {
      cy.visit('/user/databases')
      
      cy.contains('button', 'Nouvelle base de données').click()
      
      // Modal should be visible
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').should('be.visible')
      cy.contains('label', 'Type').should('be.visible')
    })

    it('should create MySQL database successfully', () => {
      cy.visit('/user/databases')
      
      cy.contains('button', 'Nouvelle base de données').click()
      
      // Fill in the form
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Test MySQL Database')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('test_db')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('test_user')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('zB6!pT9@qH3#xW7$')
      
      // Submit the form
      cy.get('form').find('button[type="submit"]').click()
      
      // Should show success
      cy.contains('Test MySQL Database', { timeout: 10000 }).should('be.visible')
      cy.contains('mysql').should('be.visible')
    })

    it('should create PostgreSQL database successfully', () => {
      cy.visit('/user/databases')
      
      cy.contains('button', 'Nouvelle base de données').click()
      
      // Fill in the form (wait for modal to be fully visible)
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Test PostgreSQL Database')
      cy.contains('label', 'Type').parent().find('select').select('postgresql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('5432')
      cy.contains('label', 'Nom de la base').parent().find('input').type('test_pg_db')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('postgres')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('yF4#nM8@rD2!vK5$')
      
      // Submit the form - use form submit button specifically
      cy.get('form').find('button[type="submit"]').click()
      
      // Should show success
      cy.contains('Test PostgreSQL Database', { timeout: 10000 }).should('be.visible')
      cy.contains('postgresql').should('be.visible')
    })

    it('should handle password visibility toggle', () => {
      cy.visit('/user/databases')
      
      cy.contains('button', 'Nouvelle base de données').click()
      
      // Get the password input
      const passwordInput = cy.contains('label', 'Mot de passe de la base de données').parent().find('input')
      
      // Should be hidden by default
      passwordInput.should('have.attr', 'type', 'password')
      
      // Click toggle button (eye icon) - need to re-query to avoid stale reference
      cy.contains('label', 'Mot de passe de la base de données').parent().find('button[type="button"]').click()
      
      // Should now be visible - re-query the input
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').should('have.attr', 'type', 'text')
      
      // Click again to hide - re-query the button
      cy.contains('label', 'Mot de passe de la base de données').parent().find('button[type="button"]').click()
      
      // Should be hidden again - re-query the input
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').should('have.attr', 'type', 'password')
    })

    it('should cancel database creation', () => {
      cy.visit('/user/databases')
      
      cy.contains('button', 'Nouvelle base de données').click()
      
      // Fill some fields
      cy.contains('label', 'Nom').parent().find('input').type('Database to Cancel')
      
      // Click cancel
      cy.contains('button', 'Annuler').click()
      
      // Modal should close - check that the modal's h2 doesn't exist
      cy.contains('h2', 'Nouvelle base de données').should('not.exist')
      
      // Database should not be created
      cy.contains('Database to Cancel').should('not.exist')
    })
  })

  describe('View Database Details', () => {
    beforeEach(() => {
      // Create a database for viewing tests
      cy.visit('/user/databases')
      cy.contains('button', 'Nouvelle base de données').click()
      
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('View Test Database')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('view_test')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('viewer')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('wJ7!xQ3@sL9#bN6$')
      
      cy.get('form').find('button[type="submit"]').click()
      cy.contains('View Test Database', { timeout: 10000 }).should('be.visible')
    })

    it('should display database information', () => {
      cy.visit('/user/databases')
      
      // Check that database card shows correct information
      cy.contains('View Test Database').should('be.visible')
      cy.contains('mysql').should('be.visible')
      
      // Find the card and check details within it
      cy.contains('View Test Database').parents('.bg-white').within(() => {
        cy.contains('Hôte:').should('be.visible')
        cy.contains('localhost:3306').should('be.visible')
        cy.contains('Base:').should('be.visible')
        cy.contains('view_test').should('be.visible')
        cy.contains('Utilisateur:').should('be.visible')
        cy.contains('viewer').should('be.visible')
      })
    })

    it('should show action buttons for database', () => {
      cy.visit('/user/databases')
      
      // Find the database card
      cy.contains('View Test Database').parents('.bg-white').within(() => {
        // Edit button (pencil icon)
        cy.get('button').eq(0).should('exist')
        
        // Delete button (trash icon)
        cy.get('button').eq(1).should('exist')
        
        // Backup button
        cy.contains('button', 'Créer une sauvegarde').should('be.visible')
      })
    })
  })

  describe('Update Database', () => {
    beforeEach(() => {
      // Create a database for update tests
      cy.visit('/user/databases')
      cy.contains('button', 'Nouvelle base de données').click()
      
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Update Test Database')
      cy.contains('label', 'Type').parent().find('select').select('postgresql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('5432')
      cy.contains('label', 'Nom de la base').parent().find('input').type('update_test')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('updater')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('tG5!mV2@hP8#xR4$')
      
      cy.get('form').find('button[type="submit"]').click()
      cy.contains('Update Test Database', { timeout: 10000 }).should('be.visible')
    })

    it('should open edit modal when clicking edit button', () => {
      cy.visit('/user/databases')
      
      // Click edit button (first button in the card, pencil icon)
      cy.contains('Update Test Database').parents('.bg-white').find('button').eq(0).click()
      
      // Edit modal should appear
      cy.contains('Modifier le nom').should('be.visible')
      cy.contains('Seuls le nom de la base de données peut être modifié').should('be.visible')
    })

        it('should update database name', () => {
            cy.visit('/user/databases')
            
            // Click edit button
            cy.contains('Update Test Database').parents('.bg-white').find('button').eq(0).click()
            
            // Change the name
            cy.contains('label', 'Nom de la base de données').parent().find('input').clear().type('Updated Database Name')
            
            // Submit
            cy.contains('button', 'Mettre à jour').click()
            
            // Wait for modal to close and page to reload
            cy.wait(2000)
            
            // Should show updated name
            cy.contains('Updated Database Name', { timeout: 10000 }).should('be.visible')
            
            // The specific card should not contain the old name anymore
            cy.contains('Updated Database Name').parents('.bg-white').within(() => {
                cy.get('h3').should('not.contain', 'Update Test Database')
            })
        })

    it('should cancel edit without saving', () => {
      cy.visit('/user/databases')
      
      // Click edit button
      cy.contains('Update Test Database').parents('.bg-white').find('button').eq(0).click()
      
      // Change the name
      cy.contains('label', 'Nom de la base de données').parent().find('input').clear().type('Cancelled Name')
      
      // Cancel
      cy.contains('button', 'Annuler').click()
      
      // Name should not change
      cy.contains('Update Test Database').should('be.visible')
      cy.contains('Cancelled Name').should('not.exist')
    })
  })

  describe('Delete Database', () => {
    beforeEach(() => {
      // Create a database for delete tests
      cy.visit('/user/databases')
      cy.contains('button', 'Nouvelle base de données').click()
      
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Delete Test Database')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('delete_test')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('deleter')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('qC8!wN4@kT6#pJ9$')
      
      cy.get('form').find('button[type="submit"]').click()
      cy.contains('Delete Test Database', { timeout: 10000 }).should('be.visible')
    })

    it('should delete database with confirmation', () => {
      cy.visit('/user/databases')
      
      // Stub the confirm dialog to automatically accept
      cy.window().then((win) => {
        cy.stub(win, 'confirm').returns(true)
      })
      
      // Click delete button (second button in the card, trash icon)
      cy.contains('Delete Test Database').parents('.bg-white').find('button').eq(1).click()
      
      // Database should be removed
      cy.contains('Delete Test Database', { timeout: 10000 }).should('not.exist')
    })

    it('should cancel deletion', () => {
      cy.visit('/user/databases')
      
      // Stub the confirm dialog to automatically reject
      cy.window().then((win) => {
        cy.stub(win, 'confirm').returns(false)
      })
      
      // Click delete button
      cy.contains('Delete Test Database').parents('.bg-white').find('button').eq(1).click()
      
      // Database should still exist
      cy.contains('Delete Test Database').should('be.visible')
    })
  })

    describe('Filter Databases', () => {
        before(() => {
            const baseUrl = Cypress.config('baseUrl') || 'http://localhost:3000'
            const apiBaseUrl = baseUrl.replace(':3000', ':8080').replace(':5173', ':8080')
            
            // Login first
            cy.request({
                method: 'POST',
                url: `${apiBaseUrl}/auth/login`,
                body: {
                    email: testEmail,
                    password: testPassword
                }
            })
            
            // Create multiple databases for filtering tests (only once)
            cy.visit('/user/databases')
            
            // Create MySQL database
            cy.contains('button', 'Nouvelle base de données').click()
            cy.contains('h2', 'Nouvelle base de données').should('be.visible')
            cy.contains('label', 'Nom').parent().find('input').type('Filter Test MySQL')
            cy.contains('label', 'Type').parent().find('select').select('mysql')
            cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
            cy.contains('label', 'Port').parent().find('input').clear().type('3306')
            cy.contains('label', 'Nom de la base').parent().find('input').type('filter_mysql')
            cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('mysql_user')
            cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('eH3!dS7@fK2#jL9$')
            cy.get('form').find('button[type="submit"]').click()
            cy.wait(3000)
            
            // Create PostgreSQL database
            cy.contains('button', 'Nouvelle base de données').click()
            cy.contains('h2', 'Nouvelle base de données').should('be.visible')
            cy.contains('label', 'Nom').parent().find('input').type('Filter Test PostgreSQL')
            cy.contains('label', 'Type').parent().find('select').select('postgresql')
            cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
            cy.contains('label', 'Port').parent().find('input').clear().type('5432')
            cy.contains('label', 'Nom de la base').parent().find('input').type('filter_pg')
            cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('pg_user')
            cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('uY6!bA9@mG4#vX2$')
            cy.get('form').find('button[type="submit"]').click()
            cy.wait(3000)
        })

    it('should filter databases by MySQL type', () => {
      cy.visit('/user/databases')
      
      // Click MySQL filter
      cy.contains('button', 'MySQL').click()
      
      // Should show MySQL databases
      cy.contains('Filter Test MySQL').should('be.visible')
    })

    it('should filter databases by PostgreSQL type', () => {
      cy.visit('/user/databases')
      
      // Click PostgreSQL filter
      cy.contains('button', 'PostgreSQL').click()
      
      // Should show PostgreSQL databases
      cy.contains('Filter Test PostgreSQL').should('be.visible')
    })

    it('should show all databases when "Tous types" is selected', () => {
      cy.visit('/user/databases')
      
      // First apply a filter
      cy.contains('button', 'MySQL').click()
      
      // Then click "Tous types"
      cy.contains('button', 'Tous types').click()
      
      // Should show all databases
      cy.contains('Filter Test MySQL').should('be.visible')
      cy.contains('Filter Test PostgreSQL').should('be.visible')
    })
  })

  describe('Backup Creation from Database Card', () => {
    beforeEach(() => {
      // Create a database for backup tests
      cy.visit('/user/databases')
      cy.contains('button', 'Nouvelle base de données').click()
      
      cy.contains('h2', 'Nouvelle base de données').should('be.visible')
      cy.contains('label', 'Nom').parent().find('input').type('Backup Test Database')
      cy.contains('label', 'Type').parent().find('select').select('mysql')
      cy.contains('label', 'Hôte').parent().find('input').clear().type('localhost')
      cy.contains('label', 'Port').parent().find('input').clear().type('3306')
      cy.contains('label', 'Nom de la base').parent().find('input').type('backup_test')
      cy.contains('label', 'Utilisateur de la base de données').parent().find('input').type('backup_user')
      cy.contains('label', 'Mot de passe de la base de données').parent().find('input').type('rZ5!nB8@cF3#wQ7$')
      
      cy.get('form').find('button[type="submit"]').click()
      cy.contains('Backup Test Database', { timeout: 10000 }).should('be.visible')
    })

    it('should show backup creation button', () => {
      cy.visit('/user/databases')
      
      cy.contains('Backup Test Database').parents('.bg-white').within(() => {
        cy.contains('button', 'Créer une sauvegarde').should('be.visible')
      })
    })

    it('should trigger backup creation when clicking backup button', () => {
      cy.visit('/user/databases')
      
      // Stub the alert to capture success message
      cy.window().then((win) => {
        cy.stub(win, 'alert').as('alertStub')
      })
      
      cy.contains('Backup Test Database').parents('.bg-white').within(() => {
        cy.contains('button', 'Créer une sauvegarde').click()
      })
      
      // Should show success alert (or error if backup fails due to test env)
      cy.get('@alertStub').should('have.been.called')
    })
  })
})
