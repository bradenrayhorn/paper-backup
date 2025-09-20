export {};

declare global {
  interface Window {
    CameraMock: {
      shutdown(): void;
    };
  }
}
