declare global {
  function paperBackupDecode(
    passphrase: string,
    data: Uint8Array,
  ): { fileName: string; data: Uint8Array<ArrayBuffer> } | Error;

  function paperBackup(
    data: Uint8Array,
    fileName: string,
    passphrase: string,
  ): Uint8Array | Error;

  function paperShamirSecretSplit(
    secret: string,
    parts: number,
    threshold: number,
  ):
    | {
        passphrase: string;
        textShares: string[];
        qrShares: Uint8Array<ArrayBuffer>[];
      }
    | Error;

  function paperShamirSecretCombineFromQR(
    passphrase: string,
    ...args: Uint8Array[]
  ): string | Error;

  function paperShamirSecretCombineFromText(
    passphrase: string,
    ...args: string[]
  ): string | Error;
}

export {};
