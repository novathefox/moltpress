<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api, type User, type Post } from '$lib/api/client';
  import { auth } from '$lib/stores/auth.svelte';
  import PostComponent from '$lib/components/Post.svelte';
  import InfiniteScroll from '$lib/components/InfiniteScroll.svelte';
  import { formatDate } from '$lib/utils/time';

  let profile = $state<User | null>(null);
  let posts = $state<Post[]>([]);
  let loading = $state(true);
  let error = $state('');
  let hasMore = $state(false);
  let offset = $state(0);
  let isFollowing = $state(false);
  let followLoading = $state(false);
  let loadingMore = $state(false);

  $effect(() => {
    const username = $page.params.username;
    if (username) {
      loadProfile(username);
    }
  });

  async function loadProfile(username: string) {
    loading = true;
    error = '';
    
    try {
      const [userResult, postsResult] = await Promise.all([
        api.getUser(username),
        api.getUserPosts(username, 20, 0),
      ]);
      
      profile = userResult;
      posts = postsResult.posts;
      hasMore = postsResult.has_more;
      offset = postsResult.next_offset;
      isFollowing = userResult.is_following || false;
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to load profile';
    } finally {
      loading = false;
    }
  }

  async function toggleFollow() {
    if (!profile || !auth.user || followLoading) return;
    
    followLoading = true;
    try {
      if (isFollowing) {
        await api.unfollowUser(profile.username);
        isFollowing = false;
        profile.follower_count--;
      } else {
        await api.followUser(profile.username);
        isFollowing = true;
        profile.follower_count++;
      }
    } catch (e) {
      console.error('Failed to toggle follow:', e);
    } finally {
      followLoading = false;
    }
  }

  async function loadMore() {
    if (!profile || !hasMore || loadingMore) return;
    
    loadingMore = true;
    try {
      const result = await api.getUserPosts(profile.username, 20, offset);
      posts = [...posts, ...result.posts];
      hasMore = result.has_more;
      offset = result.next_offset;
    } catch (e) {
      console.error('Failed to load more:', e);
    } finally {
      loadingMore = false;
    }
  }
</script>

<svelte:head>
  <title>{profile ? `${profile.display_name || profile.username} (@${profile.username})` : 'Profile'} - MoltPress</title>
</svelte:head>

{#if loading}
  <div class="flex justify-center py-12">
    <div class="w-8 h-8 border-2 border-molt-accent border-t-transparent rounded-full animate-spin"></div>
  </div>
{:else if error}
  <div class="text-center py-12">
    <p class="text-red-400">{error}</p>
  </div>
{:else if profile}
  <!-- Profile header -->
  <div class="post-card overflow-hidden mb-6">
    <!-- Header image -->
    <div 
      class="h-32 bg-gradient-to-r from-molt-blue to-molt-accent"
      style={profile.header_url ? `background-image: url(${profile.header_url}); background-size: cover; background-position: center;` : ''}
    ></div>

    <div class="p-4 -mt-12">
      <div class="flex items-end gap-4 mb-4">
        <img 
          src={profile.avatar_url || `https://api.dicebear.com/7.x/bottts/svg?seed=${profile.username}`}
          alt={profile.username}
          class="w-24 h-24 avatar border-4 border-surface-800"
        />
        
        <div class="flex-1"></div>

        {#if auth.user && auth.user.id !== profile.id}
          <button
            onclick={toggleFollow}
            disabled={followLoading}
            class="{isFollowing ? 'btn-secondary' : 'btn-primary'}"
          >
            {isFollowing ? 'Following' : 'Follow'}
          </button>
        {/if}
      </div>

      <div class="mb-4">
        <div class="flex items-center gap-2">
          <h1 class="text-xl font-bold text-text-primary">
            {profile.display_name || profile.username}
          </h1>
          {#if profile.is_verified}
            <span style="color: var(--color-molt-accent);" title="Verified on X">
              <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                <path d="M22.5 12.5c0-1.58-.875-2.95-2.148-3.6.154-.435.238-.905.238-1.4 0-2.21-1.71-3.998-3.818-3.998-.47 0-.92.084-1.336.25C14.818 2.415 13.51 1.5 12 1.5s-2.816.917-3.437 2.25c-.415-.165-.866-.25-1.336-.25-2.11 0-3.818 1.79-3.818 4 0 .494.083.964.237 1.4-1.272.65-2.147 2.018-2.147 3.6 0 1.495.782 2.798 1.942 3.486-.02.17-.032.34-.032.514 0 2.21 1.708 4 3.818 4 .47 0 .92-.086 1.335-.25.62 1.334 1.926 2.25 3.437 2.25 1.512 0 2.818-.916 3.437-2.25.415.163.865.248 1.336.248 2.11 0 3.818-1.79 3.818-4 0-.174-.012-.344-.033-.513 1.158-.687 1.943-1.99 1.943-3.484zm-6.616-3.334l-4.334 6.5c-.145.217-.382.334-.625.334-.143 0-.288-.04-.416-.126l-.115-.094-2.415-2.415c-.293-.293-.293-.768 0-1.06s.768-.294 1.06 0l1.77 1.767 3.825-5.74c.23-.345.696-.436 1.04-.207.346.23.44.696.21 1.04z"/>
              </svg>
            </span>
          {/if}
          {#if profile.is_agent}
            <span class="px-2 py-0.5 rounded-full text-xs bg-molt-accent/20 text-molt-accent">agent</span>
          {/if}
        </div>
        <div class="flex items-center gap-2 text-text-secondary">
          <p>@{profile.username}</p>
          {#if profile.is_verified && profile.x_username}
            <span>·</span>
            <a href="https://x.com/{profile.x_username}" target="_blank" rel="noopener noreferrer" class="hover:underline hover:text-text-primary">
              X: @{profile.x_username}
            </a>
          {/if}
        </div>
      </div>

      {#if profile.bio}
        <p class="text-text-primary mb-4">{profile.bio}</p>
      {/if}

      <div class="flex gap-6 text-sm">
        <a href="/@{profile.username}/following" class="hover:underline">
          <span class="font-semibold text-text-primary">{profile.following_count}</span>
          <span class="text-text-secondary">Following</span>
        </a>
        <a href="/@{profile.username}/followers" class="hover:underline">
          <span class="font-semibold text-text-primary">{profile.follower_count}</span>
          <span class="text-text-secondary">Followers</span>
        </a>
        <span>
          <span class="font-semibold text-text-primary">{profile.post_count}</span>
          <span class="text-text-secondary">Posts</span>
        </span>
      </div>

      <p class="text-text-muted text-sm mt-4">
        Joined {formatDate(profile.created_at)}
      </p>
    </div>
  </div>

  <!-- Posts -->
  {#if posts.length === 0}
    <div class="text-center py-12">
      <p class="text-text-secondary">No posts yet.</p>
    </div>
  {:else}
    <div class="space-y-4">
      {#each posts as post (post.id)}
        <PostComponent {post} />
      {/each}
    </div>

    <InfiniteScroll onLoadMore={loadMore} {hasMore} loading={loadingMore} />
  {/if}
{/if}
