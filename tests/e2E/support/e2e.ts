// ***********************************************************
// Support file for Cypress E2E tests
// This file is processed and loaded automatically before test files
// ***********************************************************

// Import commands
import './commands'

// Global configuration
Cypress.on('uncaught:exception', (err, runnable) => {
  // Prevent Cypress from failing tests on uncaught exceptions
  // You can customize this to allow specific errors
  console.error('Uncaught exception:', err.message)
  return false
})

// Before each test
beforeEach(() => {
  // Clear localStorage and sessionStorage
  cy.clearLocalStorage()
  cy.clearCookies()
})

// After each test
afterEach(function() {
  // Take screenshot on failure
  if (this.currentTest?.state === 'failed') {
    const testName = this.currentTest.title.replace(/\s+/g, '_')
    cy.screenshot(`failed/${testName}`)
  }
})

// Global types for better TypeScript support
declare global {
  namespace Cypress {
    interface Chainable {
      // Custom commands will be declared here
      login(email: string, password: string): Chainable<void>
      logout(): Chainable<void>
      registerUser(userData: {
        firstname: string
        lastname: string
        email: string
        password: string
      }): Chainable<void>
      createDatabase(dbData: {
        name: string
        type: string
        host: string
        port: string
        username: string
        password: string
        db_name: string
      }): Chainable<void>
      createSchedule(scheduleData: {
        database_id: number
        name: string
        cron_expression: string
      }): Chainable<void>
      deleteAllTestData(): Chainable<void>
      checkAccessibility(): Chainable<void>
    }
  }
}

export {}

