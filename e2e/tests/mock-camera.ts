import { Page } from '@playwright/test';

type mockCameraParams = {
  name: string;
  images: string[];
};

export async function mockCamera(params: mockCameraParams, page: Page) {
  await page.evaluate(({ name, images }: mockCameraParams) => {
    class CameraMock {
      private deviceID = 'ba171c43-9331-4a88-bc76-8a8c2cf4a5e2';
      private groupID = 'b7b40271-9456-4c84-957c-8ddde6533c32';

      private name: string;
      private imageDataUrls: string[];

      private onShutdownCamera: (() => void) | undefined;

      private canvas: HTMLCanvasElement | undefined;

      constructor(name: string, imageDataUrls: string[]) {
        this.name = name;
        this.imageDataUrls = imageDataUrls;

        Object.defineProperty(navigator, 'mediaDevices', {
          writable: true,
          configurable: true,
          value: {
            getUserMedia: this.getStream,
            enumerateDevices: this.enumerateDevices,
          },
        });
      }

      shutdown = () => {
        this.onShutdownCamera?.();
      };

      enumerateDevices = (): Promise<MediaDeviceInfo[]> => {
        const infoRaw: Omit<MediaDeviceInfo, 'toJSON'> = {
          deviceId: this.deviceID,
          groupId: this.groupID,
          kind: 'videoinput',
          label: this.name,
        };
        const info: MediaDeviceInfo = {
          ...infoRaw,
          toJSON: () => JSON.stringify(infoRaw),
        };

        return Promise.resolve([info]);
      };

      getStream = async (constraints: MediaStreamConstraints): Promise<MediaStream> => {
        this.onShutdownCamera?.();
        if (constraints.audio) return Promise.reject('audio track not supported');

        if (constraints.video === false) {
          return Promise.reject('only supports video track');
        } else if (typeof constraints.video !== 'boolean') {
          const c = constraints.video;

          if (c.deviceId && c.deviceId !== this.deviceID)
            return Promise.reject(`unknown device: ${c.deviceId}`);

          if (c.groupId && c.groupId !== this.groupID)
            return Promise.reject(`unknown groupId: ${c.groupId}`);
        }

        if (!this.canvas) {
          const images: Array<HTMLImageElement> = [];

          for (const imageDataUrl of this.imageDataUrls) {
            const image = new Image();
            image.src = imageDataUrl;
            await image.decode();
            images.push(image);
          }

          const width = images[0].naturalWidth;
          const height = images[0].naturalHeight;

          this.canvas = document.createElement('canvas');
          this.canvas.id = 'cm-canvas';

          this.canvas.width = width;
          this.canvas.height = height;

          const ctx = this.canvas.getContext('2d');

          const framesPerImage = 10;
          const maxCounter = framesPerImage * images.length;

          let counter = 0;
          const id = setInterval(() => {
            ctx.clearRect(0, 0, width, height);
            ctx.drawImage(images[Math.floor(counter / framesPerImage)], 0, 0, width, height);

            counter++;
            if (counter >= maxCounter) {
              counter = 0;
            }
          }, 1000 / 30);

          this.onShutdownCamera = () => {
            clearInterval(id);
            images.forEach((image) => {
              image.remove();
            });
            this.canvas.remove();
            this.canvas = undefined;
          };
        }

        const cv = this.canvas.captureStream(30);
        const mediaStream = new MediaStream(cv.getTracks());
        return Promise.resolve(mediaStream);
      };
    }

    const mock = new CameraMock(name, images);

    window.CameraMock = mock;
  }, params);
}

export async function closeCamera(page: Page) {
  await page.evaluate(() => {
    if (window.CameraMock) {
      window.CameraMock.shutdown();
    }
  });
}
