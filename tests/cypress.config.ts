import { defineConfig } from 'cypress'

export default defineConfig({
  e2e: {
    // Support Docker and local development
    // Docker: Frontend on port 3000, Backend on port 8080
    // Local: Frontend on port 5173, Backend on port 8080
    baseUrl: process.env.CYPRESS_BASE_URL || 'http://localhost:3000',
    specPattern: 'e2E/**/*.cy.{js,jsx,ts,tsx}',
    supportFile: 'e2E/support/e2e.ts',
    fixturesFolder: 'e2E/fixtures',
    videosFolder: 'e2E/videos',
    screenshotsFolder: 'e2E/screenshots',
    
    // Timeouts (augmentés pour Docker)
    defaultCommandTimeout: 15000,
    pageLoadTimeout: 90000,
    requestTimeout: 15000,
    responseTimeout: 45000,
    
    // Viewport
    viewportWidth: 1920,
    viewportHeight: 1080,
    
    // Video & Screenshots
    video: true,
    videoCompression: 32,
    screenshotOnRunFailure: true,
    
    // Retry
    retries: {
      runMode: 2,
      openMode: 0
    },
    
    // Environment variables
    env: {
      // Support Docker and local development
      apiUrl: process.env.CYPRESS_API_URL || 'http://localhost:8080/api',
      coverage: true,
      isDocker: process.env.CYPRESS_IS_DOCKER || 'false'
    },
    
    setupNodeEvents(on, config) {
      // implement node event listeners here
      return config
    },
  },
  
  // Component testing configuration (optional)
  component: {
    devServer: {
      framework: 'vue',
      bundler: 'vite',
    },
  },
})

