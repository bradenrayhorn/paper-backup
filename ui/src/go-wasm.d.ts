declare global {
  interface Window {
    paperBackup: (
      data: Uint8Array,
      fileName: string,
      passphrase: string,
    ) => Uint8Array;

    paperBackupDecode: (
      passphrase: string,
      data: Uint8Array,
    ) => { fileName: string; data: Uint8Array } | Error;

    paperShamirSecretSplit: (
      secret: string,
      parts: number,
      threshold: number,
    ) =>
      | { passphrase: string; textShares: string[]; qrShares: Uint8Array[] }
      | Error;

    paperShamirSecretCombineFromQR: (
      passphrase: string,
      ...args: Uint8Array[]
    ) => string | Error;

    paperShamirSecretCombineFromText: (
      passphrase: string,
      ...args: string[]
    ) => string | Error;
  }
}

export {};
