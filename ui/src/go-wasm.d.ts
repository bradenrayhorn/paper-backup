declare global {
  function paperBackupDecode(
    passphrase: string,
    data: Uint8Array,
  ): Promise<{ fileName: string; data: Uint8Array<ArrayBuffer> }>;

  function paperBackup(
    data: Uint8Array,
    fileName: string,
    passphrase: string,
  ): Promise<Uint8Array>;

  function paperShamirSecretSplit(
    secret: string,
    parts: number,
    threshold: number,
  ): Promise<{
    passphrase: string;
    textShares: string[];
    qrShares: Uint8Array<ArrayBuffer>[];
  }>;

  function paperShamirSecretCombineFromQR(
    passphrase: string,
    ...args: Uint8Array[]
  ): Promise<string>;

  function paperShamirSecretCombineFromText(
    passphrase: string,
    ...args: string[]
  ): Promise<string>;
}

export {};
