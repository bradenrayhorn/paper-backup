<script lang="ts">
  import PassphraseInput from "../lib/PassphraseInput.svelte";
  import VideoCanvas from "../lib/VideoCanvas.svelte";
  import { workerClient } from "../lib/worker-client.svelte";
  import { zxing } from "../zxing";

  let passphrase = $state("");
  let error = $state("");
  let secret = $state("");

  let seenCodes = $state<string[]>([]);
  let codeBytes = $state<Uint8Array[]>([]);
  let doneScanning = $state(false);

  async function sha256(uint8Array: Uint8Array) {
    const hashBuffer = await crypto.subtle.digest("SHA-256", uint8Array);
    return bufferToHex(hashBuffer);
  }

  function bufferToHex(buffer: ArrayBuffer) {
    return [...new Uint8Array(buffer)]
      .map((b) => b.toString(16).padStart(2, "0"))
      .join("");
  }

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

      sha256(r).then((hash) => {
        if (!seenCodes.includes(hash)) {
          seenCodes.push(hash);
          codeBytes.push(r);
        }
      });
    }
  }

  function onDoneScanning() {
    doneScanning = true;
  }

  async function onReconstruct() {
    const shares = $state.snapshot(codeBytes);

    const result = await workerClient.send(
      "shamirSecretCombineFromQR",
      {
        passphrase,
        shares,
      },
      shares.map((a) => a.buffer),
    );
    if (result instanceof Error) {
      error = result.message;
      return;
    }

    error = "";
    secret = result;
  }
</script>

<div class="wrapper">
  {#if error}
    <div>{error}</div>
  {/if}
  Combine a shamir secret. Please scan the QR codes.

  {#if doneScanning}
    <div>
      <div>Enter the secret passphrase.</div>

      <label>
        Passphrase
        <PassphraseInput bind:value={passphrase} />
      </label>

      <button onclick={onReconstruct}>Reconstruct secret.</button>
    </div>

    <div class="secret">
      {secret}
    </div>
  {:else}
    <div>Scanned codes: {seenCodes.length}</div>
    <button onclick={onDoneScanning}>Done!</button>
    <VideoCanvas
      {onFrame}
      onError={(error) => {
        console.log(error);
      }}
    />
  {/if}
</div>

<style>
  .wrapper {
    padding: 1rem;
  }
  .secret {
    font-family: monospace;
    width: 100%;
    padding: 1rem;
    border: 1px solid black;
    white-space: pre-wrap;
  }
</style>
