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
  }
}

export {};
