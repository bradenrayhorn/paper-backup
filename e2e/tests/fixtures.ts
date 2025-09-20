import fs from 'fs';

export async function loadImage(file: string): Promise<string> {
  const imageBuffer = await fs.promises.readFile(fixturePath(file));
  const base64Image = imageBuffer.toString('base64');
  return `data:image/png;base64,${base64Image}`;
}

export function fixturePath(file: string): string {
  return __dirname + `/fixtures/${file}`;
}
