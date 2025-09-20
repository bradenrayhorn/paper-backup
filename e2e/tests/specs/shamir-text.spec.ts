import shamirTextFixture from '@tests/fixtures/shamir-text';
import { expect, test } from '@tests/test';

test('can split secret and recombine', async ({ page }) => {
  await page.goto('');

  // first split the secret
  await page.getByRole('link', { name: /^Shamir Secret Split$/ }).click();
  await page.getByLabel('Secret').fill('Legumes!');
  await page.getByRole('button', { name: 'Split' }).click();

  const share1 = await page.getByLabel('Text share 1').textContent();
  const share2 = await page.getByLabel('Text share 2').textContent();
  const share3 = await page.getByLabel('Text share 3').textContent();

  await page.getByRole('button', { name: /Show encryption key/i }).click();
  const passphrase = await page.getByLabel('Passphrase').textContent();

  // now combine the secret
  await page.goto('');

  await page.getByRole('link', { name: /^Shamir Secret Combine From Text$/ }).click();

  await page.getByLabel('Secret codes').fill(`${share1},${share2},${share3}`);
  await page.getByLabel('Passphrase').fill(passphrase);

  await page.getByRole('button', { name: 'Reconstruct secret.' }).click();

  await expect(page.getByText('Legumes!')).toBeVisible();
});

test('can combine codes from fixture', async ({ page }) => {
  await page.goto('');

  await page.getByRole('link', { name: /^Shamir Secret Combine From Text$/ }).click();

  await page.getByLabel('Secret codes').fill(shamirTextFixture);
  await page.getByLabel('Passphrase').fill('A3EH KNK8 N8H9 W68A 216E');

  await page.getByRole('button', { name: 'Reconstruct secret.' }).click();

  await expect(page.getByText('Awesome secret')).toBeVisible();
});
