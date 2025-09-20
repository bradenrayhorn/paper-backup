import { fixturePath, loadImage } from '@tests/fixtures';
import { mockCamera } from '@tests/mock-camera';
import { expect, test } from '@tests/test';
import { text } from 'stream/consumers';

test('can create and restore file backup', async ({ page }) => {
  // first create the backup
  await page.goto('');

  await page.getByRole('link', { name: 'Create File Backup' }).click();
  await page.getByLabel('Enter a passphrase').fill('secret123');
  await page.getByLabel('Upload a file').setInputFiles(fixturePath('to-backup.txt'));

  await page.getByRole('button', { name: 'Backup' }).click();

  const backupBuffer = (await page.locator('svg').screenshot({ type: 'png' })).toString('base64');
  const imageData = `data:image/png;base64,${backupBuffer}`;

  // now restore the backup
  await page.goto('');

  await mockCamera({ name: 'Back Camera', images: [imageData] }, page);

  await page.getByRole('link', { name: 'Restore File Backup' }).click();
  await page.getByRole('combobox').selectOption('Back Camera');

  await expect(page.getByText('QR Code successfully scanned!')).toBeVisible();

  await page.getByLabel('Passphrase').fill('secret123');

  const downloadPromise = page.waitForEvent('download', { timeout: 120 * 1000 });
  await page.getByRole('button', { name: 'Download file' }).click();
  const result = await text(await (await downloadPromise).createReadStream());

  expect(result).toBe('Apples!\n');
});

test('can restore file backup from fixture', async ({ page }) => {
  await page.goto('');

  const imageData = await loadImage('file-backup.png');
  await mockCamera({ name: 'Back Camera', images: [imageData] }, page);

  await page.getByRole('link', { name: 'Restore File Backup' }).click();
  await page.getByRole('combobox').selectOption('Back Camera');

  await expect(page.getByText('QR Code successfully scanned!')).toBeVisible();

  await page.getByLabel('Passphrase').fill('password');

  const downloadPromise = page.waitForEvent('download', { timeout: 120 * 1000 });
  await page.getByRole('button', { name: 'Download file' }).click();
  const result = await text(await (await downloadPromise).createReadStream());

  expect(result).toBe('This is a test.\n');
});
