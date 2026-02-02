<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';

  let tags = $state<{tag: string; count: number}[]>([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      const response = await api.getTrendingTags(10);
      tags = response.tags;
    } catch (e) {
      console.error('Failed to load trending tags:', e);
      tags = [];
    } finally {
      loading = false;
    }
  });
</script>

<div class="flex items-center gap-2 overflow-x-auto scrollbar-hide">
  {#if loading}
    {#each Array(6) as _, i}
      <div class="h-8 w-20 rounded-full animate-pulse flex-shrink-0" style="background: var(--color-surface-300);"></div>
    {/each}
  {:else if tags.length > 0}
    {#each tags as { tag }}
      <a 
        href="/tagged/{tag}"
        class="tag-pill flex-shrink-0"
      >
        #{tag}
      </a>
    {/each}
  {:else}
    <span class="text-sm" style="color: var(--color-text-muted);">No trending tags yet</span>
  {/if}
</div>

<style>
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
</style>
