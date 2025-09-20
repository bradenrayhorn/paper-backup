<script lang="ts">
  import VideoCanvas from "../lib/VideoCanvas.svelte";
  import { workerClient } from "../lib/worker-client.svelte";
  import { zxing } from "../zxing";

  let passphrase = $state("");
  let error = $state("");
  let foundCode = $state<Uint8Array | null>(null);

  function onFrame(frame: Uint8ClampedArray, width: number, height: number) {
    var buffer = zxing._malloc(frame.byteLength);
    zxing.HEAPU8.set(frame, buffer);
    var result = zxing.readBarcodeFromPixmap(
      buffer,
      width,
      height,
      true,
      "QRCode",
    );
    zxing._free(buffer);

    if (result.bytes) {
      const r = new Uint8Array(result.bytes.length);
      r.set(result.bytes);
      foundCode = r;
    }
  }

  async function onDownload() {
    if (!foundCode) return;
    const result = await workerClient.send(
      "backupDecrypt",
      {
        passphrase,
        data: foundCode,
      },
      [foundCode.buffer],
    );
    if (result instanceof Error) {
      error = result.message;
      return;
    }

    error = "";

    const blob = new Blob([result.data], {});
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.download = result.fileName;
    a.href = url;
    a.click();
    a.remove();

    URL.revokeObjectURL(url);
  }
</script>

<div>
  {#if error}
    <div>{error}</div>
  {/if}
  Decrypt a file backup. Please scan the QR code backup you'd like to view.

  {#if foundCode}
    <div>
      <div>QR Code successfully scanned!</div>

      <label>
        Passphrase
        <input bind:value={passphrase} type="password" />
      </label>

      <button onclick={onDownload}>Download file</button>
    </div>
  {:else}
    <VideoCanvas
      {onFrame}
      onError={(error) => {
        console.log(error);
      }}
    />
  {/if}
</div>
