<script lang="ts">
  import type { User, Post } from '$lib/api/client';
  import ProfileContent from './ProfileContent.svelte';
  import ProfileSidebar from './ProfileSidebar.svelte';

  let {
    user,
    posts,
    loadingMore = false,
    hasMore = false,
    onLoadMore,
    onFollow,
    onUnfollow,
  }: {
    user: User;
    posts: Post[];
    loadingMore?: boolean;
    hasMore?: boolean;
    onLoadMore?: () => void;
    onFollow?: () => Promise<void>;
    onUnfollow?: () => Promise<void>;
  } = $props();

  const profileBgStyle = $derived(
    user.theme_settings?.colors?.background 
      ? `background-color: ${user.theme_settings.colors.background}` 
      : 'background-color: var(--color-surface-50)'
  );
</script>

<div class="profile-view">
  <div class="flex h-full items-start gap-4 md:gap-6">
    <div 
      class="flex-1 min-w-0 rounded-2xl overflow-hidden shadow-sm"
      style={profileBgStyle}
    >
      <ProfileContent 
        {user}
        {posts}
        {loadingMore}
        {hasMore}
        {onLoadMore}
        {onFollow}
        {onUnfollow}
      />
    </div>
    <div class="hidden md:block w-[320px] flex-shrink-0">
      <ProfileSidebar {user} />
    </div>
  </div>
</div>

<style>
  .profile-view {
    max-width: 1000px;
    height: 100%;
    margin: 0 auto;
  }
</style>
