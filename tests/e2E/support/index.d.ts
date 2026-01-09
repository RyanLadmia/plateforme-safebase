/// <reference types="cypress" />

declare namespace Cypress {
  interface Chainable {
    /**
     * Custom command to authenticate user via API
     * @example cy.authenticateUser('user@example.com', 'password')
     */
    authenticateUser(email?: string, password?: string): Chainable<void>

    /**
     * Custom command to login via UI
     * @example cy.login('user@example.com', 'password')
     */
    login(email: string, password: string): Chainable<void>

    /**
     * Custom command to logout current user
     * @example cy.logout()
     */
    logout(): Chainable<void>

    /**
     * Custom command to register a new user via UI
     * @example cy.registerUser({ firstname: 'John', lastname: 'Doe', email: 'john@example.com', password: 'password' })
     */
    registerUser(userData: {
      firstname: string
      lastname: string
      email: string
      password: string
      confirm_password?: string
    }): Chainable<void>

    /**
     * Custom command to check basic accessibility
     * @example cy.checkAccessibility()
     */
    checkAccessibility(): Chainable<void>

    /**
     * Custom command to clean up test users from database
     * @example cy.cleanupTestUsers()
     */
    cleanupTestUsers(): Chainable<void>
  }
}

