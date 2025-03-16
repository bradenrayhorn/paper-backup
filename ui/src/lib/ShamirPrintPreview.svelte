<script lang="ts">
  let {
    passphrase,
    qrShares,
    textShares,
    onClose,
  }: {
    passphrase: string;
    qrShares: string[];
    textShares: string[];
    onClose: () => void;
  } = $props();

  let qrSize = $state(100);
  let isShowingKey = $state(false);
</script>

<div>
  {#if isShowingKey}
    <button onclick={() => (isShowingKey = false)}>Show Printout</button>

    <p>
      <b>Passphrase:</b> <span class="mono">{passphrase}</span>
    </p>

    <br />

    <p>Important! Write the passphrase down on each printed share.</p>
  {:else}
    <div class="menu hide-in-print">
      <button onclick={onClose}>Go Back</button>

      <input
        type="range"
        min="50"
        max="1000"
        bind:value={qrSize}
        style="width: 100%"
      />

      <button onclick={() => (isShowingKey = true)}>SHOW ENCRYPTION KEY</button>

      <p class="wrap">
        <b>Important!</b>
        Print this page, then write the encryption key on each printed share.
      </p>
    </div>

    <div
      class="code-list"
      style="grid-template-rows: repeat({textShares.length}, 1fr) 0 2rem repeat({textShares.length}, 1fr)"
    >
      {#each textShares as share, i (i)}
        <div class="share">
          {share}
        </div>
      {/each}

      <div class="pagebreak"></div>

      <!-- Spacer given height in grid-template-rows style above. -->
      <div></div>

      {#each qrShares as svg, i (i)}
        <div class="code">
          <div style:width={`${qrSize}px`}>
            <!-- eslint-disable svelte/no-at-html-tags -->
            {@html svg}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .share {
    white-space: pre;
    font-family: monospace;
    font-size: 0.5rem;

    width: fit-content;
    padding: 0.5rem 0;
  }

  .pagebreak {
    break-before: page;
    page-break-before: always;
    height: 0;
  }

  .code-list {
    display: grid;
    grid-auto-rows: 1fr;
    padding: 0 3rem 0 3rem;
    padding-top: 2rem;
  }

  .code-list > div {
    display: flex;
    align-items: center;
  }

  .code {
    text-align: right;
    justify-content: flex-end;
    padding: 0.5rem 0;
  }

  .menu {
    margin: 1rem 0 1rem 0;
  }

  .wrap {
    margin: 1rem;
  }

  button {
    margin-left: 1rem;
  }

  .mono {
    white-space: pre;
    font-family: monospace;
  }

  @media print {
    .hide-in-print {
      display: none;
    }
  }
</style>
