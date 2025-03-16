<script lang="ts">
  import { zxing } from "../zxing";

  let files = $state<FileList | null>(null);
  let passphrase = $state("");

  let size = $state(200);
  let error = $state("");

  let qrCode = $state<string | null>(null);

  async function doBackup() {
    if (files?.length !== 1 || passphrase.trim().length < 1) {
      error = "Please enter a single file and a passphrase.";
      return;
    }

    const file = files[0];
    const resultGo = window.paperBackup(
      await file.bytes(),
      file.name ? file.name : "unknown.data",
      passphrase,
    );
    const backupData = new Uint8Array(resultGo.length);
    backupData.set(resultGo);

    const result = zxing.generateQRCodeFromBinary(
      backupData,
      "BINARY",
      0,
      10,
      10,
      5,
    );
    if (result.error) {
      error = result.error;
    } else {
      // TODO - sanitze-html to svg only
      qrCode = result.svg;

      result.delete();
    }
  }
</script>

{#if qrCode !== null}
  <div>
    <div class="hide-in-print">
      <input
        type="range"
        min="100"
        max="1000"
        bind:value={size}
        style="width: 100%"
      />
    </div>

    <div style:width={`${size}px`} style="padding: 2rem;">
      <!-- eslint-disable svelte/no-at-html-tags -->
      {@html qrCode}
    </div>
  </div>
{:else}
  <div>
    Create a file backup.

    <div>
      <label
        >Enter a passphrase
        <input bind:value={passphrase} name="passphrase" type="text" /></label
      >
    </div>

    <div>
      <label
        >Upload a file:
        <input bind:files name="file" type="file" /></label
      >
    </div>

    <button onclick={doBackup}>Backup</button>

    <span>{error}</span>
  </div>
{/if}
