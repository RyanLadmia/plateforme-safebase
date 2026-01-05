/**
 * E2E Tests - User Profile & Settings
 * Tests user profile management, password change, and settings
 * Coverage: ~10% of the application
 */

describe('User Profile & Settings', () => {
  before(() => {
    // Authenticate once
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

  describe('View Profile', () => {
    it('should navigate to profile page', () => {
      cy.visit('/user/dashboard')
      
      // Click profile link in the dashboard
      cy.contains('Mon profil').click()
      
      // Should be on profile page
      cy.url().should('include', '/profile')
    })

    it('should display user information', () => {
      cy.visit('/user/profile')
      
      // Check that user information section is visible
      cy.contains('Informations personnelles').should('be.visible')
      
      // Check that input fields exist and are visible
      cy.contains('label', 'Prénom').should('be.visible')
      cy.contains('label', 'Nom').should('be.visible')
      cy.contains('label', 'Email').should('be.visible')
      
      // Check that fields have values (not checking specific values as they depend on current user)
      cy.contains('label', 'Prénom').parent().find('input').should('not.have.value', '').and('be.visible')
      cy.contains('label', 'Nom').parent().find('input').should('not.have.value', '').and('be.visible')
      cy.contains('label', 'Email').parent().find('input').should('not.have.value', '').and('be.visible')
    })

    it('should display user role', () => {
      cy.visit('/user/profile')
      
      // Should show role label
      cy.contains('label', 'Rôle').should('be.visible')
      
      // Should show role value (Administrateur or Utilisateur)
      cy.contains('label', 'Rôle').parent().find('input').invoke('val').should('match', /Administrateur|Utilisateur/)
    })

    it('should display informations personnelles section', () => {
      cy.visit('/user/profile')
      
      // Should show personal information section
      cy.contains('Informations personnelles').should('be.visible')
    })
  })

  describe('Update Profile', () => {
    beforeEach(() => {
      cy.visit('/user/profile')
    })

    it('should enable edit mode', () => {
      // Fields should be disabled by default
      cy.contains('label', 'Prénom').parent().find('input').should('be.disabled')
      cy.contains('label', 'Nom').parent().find('input').should('be.disabled')
      cy.contains('label', 'Email').parent().find('input').should('be.disabled')
      
      // Click edit button
      cy.contains('button', 'Modifier').click()
      
      // Fields should be enabled
      cy.contains('label', 'Prénom').parent().find('input').should('not.be.disabled')
      cy.contains('label', 'Nom').parent().find('input').should('not.be.disabled')
      cy.contains('label', 'Email').parent().find('input').should('not.be.disabled')
      
      // Should show Enregistrer and Annuler buttons
      cy.contains('button', 'Enregistrer').should('be.visible')
      cy.contains('button', 'Annuler').should('be.visible')
      
      // Modifier button should not be visible
      cy.contains('button', 'Modifier').should('not.exist')
    })

    it('should update profile information successfully', () => {
      // Intercept the profile update request
      cy.intercept('PUT', '/api/profile').as('updateProfile')
      
      // Click edit button
      cy.contains('button', 'Modifier').click()
      
      // Update profile information
      cy.contains('label', 'Prénom').parent().find('input').clear().type('TestFirstName')
      cy.contains('label', 'Nom').parent().find('input').clear().type('TestLastName')
      
      // Save changes
      cy.contains('button', 'Enregistrer').click()
      
      // Wait for API call
      cy.wait('@updateProfile').its('response.statusCode').should('eq', 200)
      
      // Success message should appear
      cy.contains(/profil mis à jour avec succès/i, { timeout: 10000 }).should('be.visible')
      
      // Fields should be disabled again
      cy.contains('label', 'Prénom').parent().find('input').should('be.disabled')
      cy.contains('label', 'Nom').parent().find('input').should('be.disabled')
      
      // Modifier button should be visible again
      cy.contains('button', 'Modifier').should('be.visible')
      
      // Verify the updated values are displayed
      cy.contains('label', 'Prénom').parent().find('input').should('have.value', 'TestFirstName')
      cy.contains('label', 'Nom').parent().find('input').should('have.value', 'TestLastName')
    })

    it('should handle duplicate email error', () => {
      // Intercept the profile update request and return error
      cy.intercept('PUT', '/api/profile', {
        statusCode: 400,
        body: {
          error: 'cet email est déjà utilisé'
        }
      }).as('updateProfileError')
      
      // Click edit button
      cy.contains('button', 'Modifier').click()
      
      // Try to update with a duplicate email
      cy.contains('label', 'Email').parent().find('input').clear().type('duplicate@test.com')
      
      // Save changes
      cy.contains('button', 'Enregistrer').click()
      
      // Wait for API call
      cy.wait('@updateProfileError')
      
      // Error message should appear
      cy.contains(/cet email est déjà utilisé/i, { timeout: 5000 }).should('be.visible')
    })

    it('should cancel edit without saving', () => {
      // Store original values
      let originalFirstName: string
      
      cy.contains('label', 'Prénom').parent().find('input').invoke('val').then((val) => {
        originalFirstName = val as string
        
        // Click edit button
        cy.contains('button', 'Modifier').click()
        
        // Modify first name
        cy.contains('label', 'Prénom').parent().find('input').clear().type('ThisShouldNotBeSaved')
        
        // Click cancel
        cy.contains('button', 'Annuler').click()
        
        // Fields should be disabled again
        cy.contains('label', 'Prénom').parent().find('input').should('be.disabled')
        
        // Should restore original value
        cy.contains('label', 'Prénom').parent().find('input').should('have.value', originalFirstName)
        
        // Should not contain the modified value
        cy.contains('label', 'Prénom').parent().find('input').should('not.have.value', 'ThisShouldNotBeSaved')
      })
    })
  })

  describe('Change Password', () => {
    beforeEach(() => {
      cy.visit('/user/profile')
    })

    it('should have password change section', () => {
      // Should show password section
      cy.contains('Changer le mot de passe').should('be.visible')
      cy.contains('label', 'Mot de passe actuel').should('be.visible')
      cy.contains('label', 'Nouveau mot de passe').should('be.visible')
      cy.contains('label', 'Confirmer le mot de passe').should('be.visible')
    })

    it('should show/hide current password', () => {
      // Find current password field
      cy.contains('label', 'Mot de passe actuel').parent().find('input').should('have.attr', 'type', 'password')
      
      // Click toggle button
      cy.contains('label', 'Mot de passe actuel').parent().find('button').click()
      
      // Should show password
      cy.contains('label', 'Mot de passe actuel').parent().find('input').should('have.attr', 'type', 'text')
    })

    it('should show/hide new password', () => {
      // Find new password field
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').should('have.attr', 'type', 'password')
      
      // Click toggle button
      cy.contains('label', 'Nouveau mot de passe').parent().find('button').click()
      
      // Should show password
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').should('have.attr', 'type', 'text')
    })

    it('should validate password fields are required', () => {
      // Intercept the password change request
      cy.intercept('PUT', '/api/profile/password').as('changePassword')
      
      // Click change password button without filling fields
      cy.contains('button', 'Changer le mot de passe').click()
      
      // Error message should appear
      cy.contains(/tous les champs sont requis/i, { timeout: 5000 }).should('be.visible')
      
      // API should not be called
      cy.get('@changePassword.all').should('have.length', 0)
    })

    it('should validate passwords match', () => {
      // Intercept the password change request
      cy.intercept('PUT', '/api/profile/password').as('changePassword')
      
      // Fill password fields with non-matching passwords
      cy.contains('label', 'Mot de passe actuel').parent().find('input').type('xK9#mQ2$vL7@wP4!nR8%')
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').type('NewP@ssw0rd123')
      cy.contains('label', 'Confirmer le mot de passe').parent().find('input').type('DifferentP@ssw0rd')
      
      // Click change password button
      cy.contains('button', 'Changer le mot de passe').click()
      
      // Error message should appear
      cy.contains(/les mots de passe ne correspondent pas/i, { timeout: 5000 }).should('be.visible')
      
      // API should not be called
      cy.get('@changePassword.all').should('have.length', 0)
    })

    it('should validate password minimum length', () => {
      // Intercept the password change request
      cy.intercept('PUT', '/api/profile/password').as('changePassword')
      
      // Fill password fields with short password
      cy.contains('label', 'Mot de passe actuel').parent().find('input').type('xK9#mQ2$vL7@wP4!nR8%')
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').type('Short1!')
      cy.contains('label', 'Confirmer le mot de passe').parent().find('input').type('Short1!')
      
      // Click change password button
      cy.contains('button', 'Changer le mot de passe').click()
      
      // Error message should appear
      cy.contains(/le mot de passe doit contenir au moins 8 caractères/i, { timeout: 5000 }).should('be.visible')
      
      // API should not be called
      cy.get('@changePassword.all').should('have.length', 0)
    })

    it('should handle incorrect current password error', () => {
      // Intercept the password change request and return error
      cy.intercept('PUT', '/api/profile/password', {
        statusCode: 400,
        body: {
          error: 'mot de passe actuel incorrect'
        }
      }).as('changePasswordError')
      
      // Fill password fields
      cy.contains('label', 'Mot de passe actuel').parent().find('input').type('WrongPassword123')
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').type('NewP@ssw0rd123456')
      cy.contains('label', 'Confirmer le mot de passe').parent().find('input').type('NewP@ssw0rd123456')
      
      // Click change password button
      cy.contains('button', 'Changer le mot de passe').click()
      
      // Wait for API call
      cy.wait('@changePasswordError')
      
      // Error message should appear
      cy.contains(/mot de passe actuel incorrect/i, { timeout: 5000 }).should('be.visible')
    })

    it('should change password successfully', () => {
      // Intercept the password change request
      cy.intercept('PUT', '/api/profile/password', {
        statusCode: 200,
        body: {
          message: 'Mot de passe changé avec succès'
        }
      }).as('changePasswordSuccess')
      
      // Fill password fields
      cy.contains('label', 'Mot de passe actuel').parent().find('input').type('xK9#mQ2$vL7@wP4!nR8%')
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').type('NewP@ssw0rd123456')
      cy.contains('label', 'Confirmer le mot de passe').parent().find('input').type('NewP@ssw0rd123456')
      
      // Click change password button
      cy.contains('button', 'Changer le mot de passe').click()
      
      // Wait for API call
      cy.wait('@changePasswordSuccess')
      
      // Success message should appear
      cy.contains(/mot de passe changé avec succès/i, { timeout: 10000 }).should('be.visible')
      
      // Form should be cleared
      cy.contains('label', 'Mot de passe actuel').parent().find('input').should('have.value', '')
      cy.contains('label', 'Nouveau mot de passe').parent().find('input').should('have.value', '')
      cy.contains('label', 'Confirmer le mot de passe').parent().find('input').should('have.value', '')
    })
  })

  describe('User Statistics', () => {
    beforeEach(() => {
      cy.visit('/user/profile')
    })

    it('should display statistics section', () => {
      // Should show statistics section
      cy.contains('Statistiques du compte').should('be.visible')
    })

    it('should display database count', () => {
      // Should show number of databases
      cy.contains('Bases de données').should('be.visible')
      cy.get('.text-3xl.font-bold.text-blue-600').should('exist')
    })

    it('should display backup count', () => {
      // Should show number of backups
      cy.contains('Sauvegardes totales').should('be.visible')
      cy.get('.text-3xl.font-bold.text-green-600').should('exist')
    })

    it('should display completed backups count', () => {
      // Should show number of completed backups
      cy.contains('Sauvegardes réussies').should('be.visible')
      cy.get('.text-3xl.font-bold.text-orange-600').should('exist')
    })
  })
})

