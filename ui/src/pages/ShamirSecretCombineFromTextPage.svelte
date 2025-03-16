<script lang="ts">
  import PassphraseInput from "../lib/PassphraseInput.svelte";

  let passphrase = $state("");
  let error = $state("");
  let secret = $state("");

  let codeInput = $state<string>("");

  function onReconstruct() {
    const codes = codeInput
      .toUpperCase()
      .trim()
      .split(",")
      .map((c) => c.trim())
      .filter((c) => c.length > 0);
    const result = window.paperShamirSecretCombineFromText(
      passphrase,
      ...codes,
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
  Combine a shamir secret. Please enter codes separated by a comma.

  <div>
    <div>Enter the secret passphrase.</div>

    <label>
      Passphrase
      <PassphraseInput bind:value={passphrase} />
    </label>

    <button onclick={onReconstruct}>Reconstruct secret.</button>
  </div>

  {#if secret}
    <div class="secret">
      {secret}
    </div>
  {:else}
    <textarea
      bind:value={codeInput}
      autocomplete="off"
      autocapitalize="none"
      spellcheck="false"
    ></textarea>
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
