import { loadImage } from '@tests/fixtures';
import { mockCamera } from '@tests/mock-camera';
import { expect, test } from '@tests/test';

test('can split secret and recombine', async ({ page }) => {
  await page.goto('');

  // first split the secret
  await page.getByRole('link', { name: /^Shamir Secret Split$/ }).click();
  await page.getByLabel('Secret').fill('Broccoli!');
  await page.getByRole('button', { name: 'Split' }).click();

  const buffer1 = (await page.locator('svg').nth(0).screenshot({ type: 'png' })).toString('base64');
  const buffer2 = (await page.locator('svg').nth(1).screenshot({ type: 'png' })).toString('base64');
  const buffer3 = (await page.locator('svg').nth(2).screenshot({ type: 'png' })).toString('base64');

  await page.getByRole('button', { name: /Show encryption key/i }).click();
  const passphrase = await page.getByLabel('Passphrase').textContent();

  // now combine the secret
  await page.goto('');

  await mockCamera(
    {
      name: 'Back Camera',
      images: [
        `data:image/png;base64,${buffer1}`,
        `data:image/png;base64,${buffer2}`,
        `data:image/png;base64,${buffer3}`,
      ],
    },
    page,
  );

  await page.getByRole('link', { name: /^Shamir Secret Combine$/ }).click();
  await page.getByRole('combobox').selectOption('Back Camera');

  await expect(page.getByText('Scanned codes: 3')).toBeVisible();
  await page.getByRole('button', { name: 'Done!' }).click();

  await page.getByLabel('Passphrase').fill(passphrase);

  await page.getByRole('button', { name: 'Reconstruct secret.' }).click();

  await expect(page.getByText('Broccoli!')).toBeVisible();
});

test('can combine QR codes from fixture', async ({ page }) => {
  await page.goto('');

  const imageData1 = await loadImage('shamir1.png');
  const imageData2 = await loadImage('shamir2.png');
  const imageData3 = await loadImage('shamir3.png');
  await mockCamera({ name: 'Back Camera', images: [imageData1, imageData2, imageData3] }, page);

  await page.getByRole('link', { name: /^Shamir Secret Combine$/ }).click();
  await page.getByRole('combobox').selectOption('Back Camera');

  await expect(page.getByText('Scanned codes: 3')).toBeVisible();
  await page.getByRole('button', { name: 'Done!' }).click();

  await page.getByLabel('Passphrase').fill('942931E86676HNX67NEK');

  await page.getByRole('button', { name: 'Reconstruct secret.' }).click();

  await expect(page.getByText('Big secret')).toBeVisible();
});
