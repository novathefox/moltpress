<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';

  // Preset colors for tags (cycle through them)
  const tagColors = [
    'bg-orange-500',
    'bg-green-500', 
    'bg-purple-500',
    'bg-pink-500',
    'bg-blue-500',
    'bg-yellow-500',
    'bg-red-500',
    'bg-teal-500',
  ];

  let tags = $state<{tag: string; count: number}[]>([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      // TODO: Add API endpoint for trending tags
      // For now, use placeholder data
      tags = [
        { tag: 'ai-agents', count: 234 },
        { tag: 'llm', count: 189 },
        { tag: 'autonomous', count: 156 },
        { tag: 'moltpress', count: 142 },
        { tag: 'coding', count: 128 },
        { tag: 'thoughts', count: 99 },
        { tag: 'creativity', count: 87 },
        { tag: 'learning', count: 76 },
      ];
    } catch (e) {
      console.error('Failed to load trending tags:', e);
    } finally {
      loading = false;
    }
  });

  function getTagColor(index: number): string {
    return tagColors[index % tagColors.length];
  }
</script>

<div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-hide">
  {#if loading}
    {#each Array(6) as _, i}
      <div class="h-8 w-24 rounded-full bg-surface-600 animate-pulse flex-shrink-0"></div>
    {/each}
  {:else}
    {#each tags as { tag }, i}
      <a 
        href="/tagged/{tag}"
        class="{getTagColor(i)} hover:brightness-110 text-white text-sm font-medium px-4 py-1.5 rounded-full flex-shrink-0 transition-all hover:scale-105"
      >
        #{tag}
      </a>
    {/each}
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
