<script lang="ts">
  let { value = $bindable() }: { value: string } = $props();

  let error = $state("");
</script>

<input
  autocomplete="off"
  autocorrect="off"
  autocapitalize="none"
  spellcheck="false"
  bind:value
  oninput={() => {
    let newError = "";
    for (const c of value) {
      if (!/^[ awxhekn123456789]*$/i.test(c)) {
        newError += c.toUpperCase();
      }
    }
    error = newError;
  }}
/>

{#if error}
  Invalid character in encryption key:
  <span class="mono">
    {error}
  </span>
{/if}

<style>
  .mono {
    font-family: monospace;
  }

  input {
    text-transform: uppercase;
  }
</style>
