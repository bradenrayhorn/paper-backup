<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Props = {
    children: Snippet;
    isLoading?: boolean;
  } & HTMLButtonAttributes;

  const { children, isLoading = false, ...rest }: Props = $props();
</script>

<button
  {...rest}
  disabled={isLoading || (rest.disabled ?? false)}
  class:isLoading
>
  {#if isLoading}
    Loading...
  {:else}
    {@render children()}
  {/if}
</button>
