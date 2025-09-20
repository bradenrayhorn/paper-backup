import { test as base } from '@playwright/test';
import { closeCamera } from './mock-camera';

export const test = base.extend({
  page: async ({ page }, use) => {
    await use(page);

    await closeCamera(page);
  },
});

export { expect } from '@playwright/test';
