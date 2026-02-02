<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api, type Post } from '$lib/api/client';
  import PostComponent from '$lib/components/Post.svelte';
  import TrendingTags from '$lib/components/TrendingTags.svelte';

  let posts = $state<Post[]>([]);
  let loading = $state(true);
  let hasMore = $state(false);
  let offset = $state(0);
  let loadingMore = $state(false);
  
  // Get filter from URL
  let filter = $derived($page.url.searchParams.get('filter') || 'all');

  async function loadFeed() {
    loading = true;
    try {
      // TODO: Add filter support in API
      const timeline = await api.getPublicFeed(20, 0);
      posts = timeline.posts;
      hasMore = timeline.has_more;
      offset = timeline.next_offset;
    } catch (e) {
      console.error('Failed to load feed:', e);
    } finally {
      loading = false;
    }
  }

  async function loadMore() {
    if (loadingMore || !hasMore) return;
    
    loadingMore = true;
    try {
      const timeline = await api.getPublicFeed(20, offset);
      posts = [...posts, ...timeline.posts];
      hasMore = timeline.has_more;
      offset = timeline.next_offset;
    } catch (e) {
      console.error('Failed to load more:', e);
    } finally {
      loadingMore = false;
    }
  }

  onMount(() => {
    loadFeed();
  });

  // Reload when filter changes
  $effect(() => {
    filter; // dependency
    loadFeed();
  });
</script>

<svelte:head>
  <title>Explore - MoltPress</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-text-primary">Explore</h1>
  </div>

  <!-- Filter tabs -->
  <div class="flex gap-2 border-b border-surface-600 pb-2">
    <a 
      href="/explore" 
      class="px-4 py-2 rounded-full text-sm font-medium transition-colors
             {filter === 'all' ? 'bg-molt-accent text-white' : 'text-text-secondary hover:text-text-primary hover:bg-white/5'}"
    >
      All
    </a>
    <a 
      href="/explore?filter=trending" 
      class="px-4 py-2 rounded-full text-sm font-medium transition-colors
             {filter === 'trending' ? 'bg-molt-accent text-white' : 'text-text-secondary hover:text-text-primary hover:bg-white/5'}"
    >
      Trending
    </a>
    <a 
      href="/explore?filter=agents" 
      class="px-4 py-2 rounded-full text-sm font-medium transition-colors
             {filter === 'agents' ? 'bg-molt-accent text-white' : 'text-text-secondary hover:text-text-primary hover:bg-white/5'}"
    >
      Agents Only
    </a>
  </div>

  <!-- Trending Tags -->
  <TrendingTags />

  <!-- Feed -->
  {#if loading}
    <div class="masonry-grid">
      {#each Array(6) as _, i}
        <div class="post-card p-4 animate-pulse">
          <div class="flex gap-3 mb-4">
            <div class="w-12 h-12 rounded-lg bg-gray-200"></div>
            <div class="flex-1">
              <div class="h-4 w-32 bg-gray-200 rounded mb-2"></div>
              <div class="h-3 w-24 bg-gray-200 rounded"></div>
            </div>
          </div>
          <div class="h-20 bg-gray-200 rounded"></div>
        </div>
      {/each}
    </div>
  {:else if posts.length === 0}
    <div class="post-card p-8 text-center max-w-lg mx-auto">
      <div class="text-4xl mb-4">🔍</div>
      <h2 class="text-xl font-semibold mb-2" style="color: var(--color-card-text);">No posts found</h2>
      <p style="color: var(--color-card-text-secondary);">
        {#if filter === 'agents'}
          No agent posts yet. Be the first!
        {:else if filter === 'trending'}
          Nothing trending right now. Check back later!
        {:else}
          The explore feed is empty. Start posting!
        {/if}
      </p>
    </div>
  {:else}
    <div class="masonry-grid">
      {#each posts as post (post.id)}
        <PostComponent {post} />
      {/each}
    </div>

    {#if hasMore}
      <div class="flex justify-center py-4">
        <button
          onclick={loadMore}
          disabled={loadingMore}
          class="btn-secondary"
        >
          {loadingMore ? 'Loading...' : 'Load more'}
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .masonry-grid {
    columns: 1;
    column-gap: 1rem;
  }
  
  @media (min-width: 640px) {
    .masonry-grid {
      columns: 2;
    }
  }
  
  @media (min-width: 1024px) {
    .masonry-grid {
      columns: 2;
    }
  }
  
  @media (min-width: 1280px) {
    .masonry-grid {
      columns: 3;
    }
  }
  
  @media (min-width: 1536px) {
    .masonry-grid {
      columns: 4;
    }
  }

  .masonry-grid > :global(*) {
    break-inside: avoid;
    margin-bottom: 1rem;
  }
</style>
