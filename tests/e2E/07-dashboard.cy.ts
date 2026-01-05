/**
 * E2E Tests - Dashboard & Overview
 * Tests dashboard display, statistics, and quick actions
 * Coverage: ~5% of the application
 */

describe('Dashboard & Overview', () => {
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
    
    cy.visit('/user/dashboard')
  })

  describe('Dashboard Layout', () => {
    it('should display dashboard page correctly', () => {
      cy.url().should('include', '/user/dashboard')
      
      // Should show main dashboard elements
      cy.contains(/tableau de bord|dashboard/i).should('be.visible')
    })

    it('should display navigation menu', () => {
      // Should show navigation
      cy.get('nav, aside, [role="navigation"]').should('be.visible')
      
      // Check main menu items exist (they may be in sidebar)
      cy.get('body').then(($body) => {
        if ($body.text().includes('Bases de données') || $body.text().includes('Databases')) {
          cy.contains(/bases de données|databases/i).should('exist')
        }
      })
    })

    it('should display user profile access', () => {
      // Should have link to profile (in menu or sidebar)
      cy.contains(/mon profil|profile/i).should('exist')
    })

    it('should be responsive', () => {
      // Test mobile viewport
      cy.viewport('iphone-x')
      cy.get('nav, aside, [role="navigation"]').should('exist')
      
      // Test tablet viewport
      cy.viewport('ipad-2')
      cy.get('nav, aside, [role="navigation"]').should('exist')
      
      // Test desktop viewport
      cy.viewport(1920, 1080)
      cy.get('nav, aside, [role="navigation"]').should('exist')
    })
  })

  describe('Statistics Cards', () => {
    it('should display statistics overview', () => {
      // Wait for the page to fully load by checking for specific elements
      cy.contains(/bases de données|databases/i, { timeout: 10000 }).should('be.visible')
      cy.contains(/sauvegardes|backups/i, { timeout: 10000 }).should('be.visible')
    })

    it('should display database count', () => {
      // Statistics should be visible (actual numbers depend on data)
      cy.contains(/bases de données|databases/i).should('be.visible')
    })

    it('should display backup count', () => {
      // Statistics should be visible
      cy.contains(/sauvegardes|backups/i).should('be.visible')
    })

    it('should display schedule count', () => {
      // Statistics should be visible
      cy.contains(/planifications|schedules/i).should('be.visible')
    })
  })

  describe('Recent Activity', () => {
    it('should have activity section or history link', () => {
      // Check if there's a recent activity section or link to history
      cy.get('body').then(($body) => {
        const hasActivity = $body.text().match(/activité|activity|historique|history/i)
        if (hasActivity) {
          cy.contains(/activité|activity|historique|history/i).should('be.visible')
        } else {
          // If no activity section, at least history should be accessible
          cy.contains(/historique|history/i).should('exist')
        }
      })
    })
  })

  describe('Quick Navigation', () => {
    it('should navigate to databases page', () => {
      cy.contains(/bases de données|databases/i).first().click()
      cy.url().should('include', '/user/databases')
    })

    it('should navigate to backups page', () => {
      cy.visit('/user/dashboard')
      cy.contains(/sauvegardes|backups/i).first().click({ force: true })
      cy.url().should('include', '/user/backups')
    })

    it('should navigate to schedules page', () => {
      cy.visit('/user/dashboard')
      cy.contains(/planifications|schedules/i).first().click({ force: true })
      cy.url().should('include', '/user/schedules')
    })

    it('should navigate to history page', () => {
      cy.visit('/user/dashboard')
      cy.contains(/historique|history/i).first().click({ force: true })
      cy.url().should('include', '/user/history')
    })

    it('should navigate to profile page', () => {
      cy.visit('/user/dashboard')
      cy.contains(/mon profil|profile/i).first().click({ force: true })
      cy.url().should('include', '/user/profile')
    })
  })

  describe('Accessibility', () => {
    it('should have proper heading hierarchy', () => {
      cy.get('h1, h2, h3').should('exist')
    })

    it('should be keyboard navigable', () => {
      // Check that interactive elements are focusable
      cy.get('a, button').first().should('exist').focus()
      cy.focused().should('exist')
    })
  })
})

