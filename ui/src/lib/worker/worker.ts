import wasmReady from "../../wasm/load_worker";
import type {
  MessageTypes,
  PromisedReply,
  Reply,
  Request,
  WorkerQueryMessageSet,
} from "./message-types";

self.onmessage = async (e: MessageEvent<WorkerQueryMessageSet>) => {
  await wasmReady;
  const message = e.data;

  try {
    const { data, transfer } = await handleMessage(message);
    self.postMessage({ id: message.id, data }, { transfer });
  } catch (error) {
    self.postMessage({
      id: message.id,
      data: new Error("wasm error", { cause: error }),
    });
  }
};

self.onerror = (e) => {
  console.error("error in worker", e);
};

async function handleMessage(
  message: WorkerQueryMessageSet,
): Promise<Reply<MessageTypes>> {
  switch (message.type) {
    case "backupEncrypt":
      return backupEncrypt(message.data);
    case "backupDecrypt":
      return backupDecrypt(message.data);
    case "shamirSecretSplit":
      return shamirSecretSplit(message.data);
    case "shamirSecretCombineFromQR":
      return shamirSecretCombineFromQR(message.data);
    case "shamirSecretCombineFromText":
      return shamirSecretCombineFromText(message.data);
  }
}

async function backupEncrypt({
  data,
  fileName,
  passphrase,
}: Request<"backupEncrypt">): PromisedReply<"backupEncrypt"> {
  const backup = await paperBackup(data, fileName, passphrase);

  return { data: backup, transfer: [backup.buffer] };
}

async function backupDecrypt({
  data,
  passphrase,
}: Request<"backupDecrypt">): PromisedReply<"backupDecrypt"> {
  const decrypted = await paperBackupDecode(passphrase, data);

  return { data: decrypted, transfer: [decrypted.data.buffer] };
}

async function shamirSecretSplit({
  secret,
  parts,
  threshold,
}: Request<"shamirSecretSplit">): PromisedReply<"shamirSecretSplit"> {
  const result = await paperShamirSecretSplit(secret, parts, threshold);

  return { data: result, transfer: result.qrShares.map((s) => s.buffer) };
}

async function shamirSecretCombineFromQR({
  passphrase,
  shares,
}: Request<"shamirSecretCombineFromQR">): PromisedReply<"shamirSecretCombineFromQR"> {
  const result = await paperShamirSecretCombineFromQR(passphrase, ...shares);

  return { data: result };
}

async function shamirSecretCombineFromText({
  passphrase,
  shares,
}: Request<"shamirSecretCombineFromText">): PromisedReply<"shamirSecretCombineFromText"> {
  const result = await paperShamirSecretCombineFromText(passphrase, ...shares);

  return { data: result };
}
