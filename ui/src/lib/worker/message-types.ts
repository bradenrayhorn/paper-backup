export type Messages = {
  backupEncrypt: {
    request: { data: Uint8Array; fileName: string; passphrase: string };
    reply: Uint8Array;
  };
  backupDecrypt: {
    request: { data: Uint8Array; passphrase: string };
    reply: { fileName: string; data: Uint8Array<ArrayBuffer> };
  };

  shamirSecretSplit: {
    request: { secret: string; parts: number; threshold: number };
    reply: { passphrase: string; textShares: string[]; qrShares: Uint8Array[] };
  };
  shamirSecretCombineFromQR: {
    request: { passphrase: string; shares: Uint8Array[] };
    reply: string;
  };
  shamirSecretCombineFromText: {
    request: { passphrase: string; shares: string[] };
    reply: string;
  };
};

export type MessageTypes = keyof Messages;

export type Request<T extends keyof Messages> = Messages[T]["request"];

export type Reply<T extends keyof Messages> = {
  data: Messages[T]["reply"];
  transfer?: Transferable[];
};

export type ReplyMessage = {
  id: string;
  data: Reply<keyof Messages>["data"];
};

export type WorkerQueryMessageSet = {
  [K in keyof Messages]: {
    id: string;
    type: K;
    data: Messages[K]["request"];
  };
}[keyof Messages];

export type WorkerQueryMessage<T extends MessageTypes> = {
  id: string;
  type: T;
  data: Messages[T]["request"];
};
