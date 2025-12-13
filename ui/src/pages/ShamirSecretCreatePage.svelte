<script lang="ts">
  import Button from "../lib/Button.svelte";
  import ShamirPrintPreview from "../lib/ShamirPrintPreview.svelte";
  import { workerClient } from "../lib/worker-client.svelte";
  import { zxing } from "../zxing";

  let parts = $state(5);
  let threshold = $state(3);
  let secret = $state("");

  let passphrase = $state("");
  let error = $state("");

  let isLoading = $state(false);

  let printState = $state<{
    textShares: string[];
    qrShares: string[];
  } | null>(null);

  async function doSplit() {
    const result = await workerClient.send("shamirSecretSplit", {
      secret,
      parts,
      threshold,
    });
    if (result instanceof Error) {
      error = result.message;
      return;
    }

    error = "";
    passphrase = result.passphrase;

    const newPrintState: { textShares: string[]; qrShares: string[] } = {
      textShares: result.textShares,
      qrShares: [],
    };

    for (const share of result.qrShares) {
      const result = zxing.generateQRCodeFromBinary(
        share,
        "BINARY",
        0,
        10,
        10,
        7,
      );
      if (result.error) {
        error = result.error;
      } else {
        // TODO - sanitze-html to svg only
        newPrintState.qrShares.push(result.svg);

        result.delete();
      }
    }

    printState = newPrintState;
  }
</script>

{#if printState !== null}
  <ShamirPrintPreview
    onClose={() => {
      printState = null;
    }}
    qrShares={printState.qrShares}
    textShares={printState.textShares}
    {passphrase}
  />
{:else}
  <div>
    Create a shamir secret split.

    <div>
      <label
        >Parts
        <input bind:value={parts} type="number" /></label
      >
    </div>

    <div>
      <label
        >Threshold
        <input bind:value={threshold} type="number" /></label
      >
    </div>

    <div>
      <label
        >Secret
        <textarea bind:value={secret}></textarea></label
      >
    </div>

    <Button
      {isLoading}
      onclick={() => {
        isLoading = true;
        doSplit().finally(() => {
          isLoading = false;
        });
      }}>Split</Button
    >
    <span>{error}</span>
  </div>
{/if}
