import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: undefined,
  reporter: process.env.CI ? 'html' : 'list',
  /* Firefox WASM is very slow in Playwright environment, so timeout must be increased. */
  timeout: 300 * 1000,
  expect: { timeout: 90 * 1000 },
  use: {
    baseURL: `http://localhost:8080/_index.html`,
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
      },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],

  webServer: {
    command: 'caddy file-server --listen :8080',
    url: 'http://localhost:8080/_index.html',
    reuseExistingServer: !process.env.CI,
  },
});
