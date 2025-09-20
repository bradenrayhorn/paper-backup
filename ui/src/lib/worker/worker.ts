import wasmReady from "../../wasm/load_worker";
import type {
  MessageTypes,
  Reply,
  Request,
  WorkerQueryMessageSet,
} from "./message-types";

self.onmessage = async (e: MessageEvent<WorkerQueryMessageSet>) => {
  await wasmReady;
  const message = e.data;

  try {
    const { data, transfer } = handleMessage(message);
    self.postMessage({ id: message.id, data }, { transfer });
  } catch (error) {
    self.postMessage({
      id: message.id,
      data: new Error("wasm error", { cause: error }),
    });
  }
};

function handleMessage(message: WorkerQueryMessageSet): Reply<MessageTypes> {
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

function backupEncrypt({
  data,
  fileName,
  passphrase,
}: Request<"backupEncrypt">): Reply<"backupEncrypt"> {
  const backup = paperBackup(data, fileName, passphrase);
  if (backup instanceof Error) throw backup;

  return { data: backup, transfer: [backup.buffer] };
}

function backupDecrypt({
  data,
  passphrase,
}: Request<"backupDecrypt">): Reply<"backupDecrypt"> {
  const decrypted = paperBackupDecode(passphrase, data);
  if (decrypted instanceof Error) throw decrypted;

  return { data: decrypted, transfer: [decrypted.data.buffer] };
}

function shamirSecretSplit({
  secret,
  parts,
  threshold,
}: Request<"shamirSecretSplit">): Reply<"shamirSecretSplit"> {
  const result = paperShamirSecretSplit(secret, parts, threshold);
  if (result instanceof Error) throw result;

  return { data: result, transfer: result.qrShares.map((s) => s.buffer) };
}

function shamirSecretCombineFromQR({
  passphrase,
  shares,
}: Request<"shamirSecretCombineFromQR">): Reply<"shamirSecretCombineFromQR"> {
  const result = paperShamirSecretCombineFromQR(passphrase, ...shares);
  if (result instanceof Error) throw result;

  return { data: result };
}

function shamirSecretCombineFromText({
  passphrase,
  shares,
}: Request<"shamirSecretCombineFromText">): Reply<"shamirSecretCombineFromText"> {
  const result = paperShamirSecretCombineFromText(passphrase, ...shares);
  if (result instanceof Error) throw result;

  return { data: result };
}
