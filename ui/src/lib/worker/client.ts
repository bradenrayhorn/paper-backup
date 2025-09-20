import type {
  Messages,
  MessageTypes,
  ReplyMessage,
  WorkerQueryMessage,
} from "./message-types";

type MessageRequest<T extends MessageTypes> = Messages[T]["request"];

type MessageReply<T extends MessageTypes> = Messages[T]["reply"];

type MessageQueueItem = {
  resolve: (value: MessageReply<MessageTypes> | Error) => void;
};

export class WorkerClient {
  #messageQueue: Record<string, MessageQueueItem> = {};
  #worker: Worker | undefined;

  private async init(): Promise<void> {
    if (this.#worker !== undefined) return;

    const worker_import = await import("./worker.ts?worker&inline");
    this.#worker = new worker_import.default();

    this.#worker.onmessage = this.onMessage;
  }

  private onMessage = (e: MessageEvent<ReplyMessage>) => {
    const { id, data } = e.data;

    if (!this.#messageQueue[id]) {
      console.error("unknown message from worker:", e.data);
      return;
    }

    const { resolve } = this.#messageQueue[id];
    delete this.#messageQueue[id];
    if (data instanceof Error) {
      console.error("error from worker", data);
      resolve(data);
    } else {
      resolve(data);
    }
  };

  send = async <T extends MessageTypes>(
    type: T,
    data: MessageRequest<T>,
    transfer?: Transferable[],
  ): Promise<MessageReply<T> | Error> => {
    await this.init();

    const id = (Math.random() + 1).toString(36).substring(7);

    const message: WorkerQueryMessage<T> = {
      id,
      type,
      data,
    };

    return new Promise((resolve) => {
      this.#messageQueue[id] = { resolve };
      this.#worker?.postMessage(message, { transfer });

      setTimeout(() => {
        if (this.#messageQueue[id]) {
          console.error("worker timeout", id);
          resolve(new Error("message timeout!"));
        }
      }, 90 * 1000);
    });
  };
}
