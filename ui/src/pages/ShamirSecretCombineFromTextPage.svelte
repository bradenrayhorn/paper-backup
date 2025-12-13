<script lang="ts">
  import Button from "../lib/Button.svelte";
  import PassphraseInput from "../lib/PassphraseInput.svelte";
  import { workerClient } from "../lib/worker-client.svelte";

  let passphrase = $state("");
  let error = $state("");
  let secret = $state("");
  let isLoading = $state(false);

  let codeInput = $state<string>("");

  async function onReconstruct() {
    const codes = codeInput
      .toUpperCase()
      .trim()
      .split(",")
      .map((c) => c.trim())
      .filter((c) => c.length > 0);

    const result = await workerClient.send("shamirSecretCombineFromText", {
      passphrase,
      shares: codes,
    });
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
  Combine a shamir secret. Please enter codes separated by a comma.

  <div>
    <div>Enter the secret passphrase.</div>

    <label>
      Passphrase
      <PassphraseInput bind:value={passphrase} />
    </label>

    <Button
      {isLoading}
      onclick={() => {
        isLoading = true;
        onReconstruct().finally(() => {
          isLoading = false;
        });
      }}>Reconstruct secret.</Button
    >
  </div>

  {#if secret}
    <div class="secret">
      {secret}
    </div>
  {:else}
    <label>
      Secret codes
      <textarea
        bind:value={codeInput}
        autocomplete="off"
        autocapitalize="none"
        spellcheck="false"
      ></textarea>
    </label>
  {/if}
</div>

<style>
  .wrapper {
    padding: 1rem;
  }
  .secret {
    margin-top: 1rem;
    font-family: monospace;
    width: 100%;
    padding: 1rem;
    border: 1px solid black;
    white-space: pre-wrap;
  }

  textarea {
    margin-top: 1rem;
    width: 100%;
    min-height: 400px;
    font-family: monospace;
    text-transform: uppercase;
  }
</style>
