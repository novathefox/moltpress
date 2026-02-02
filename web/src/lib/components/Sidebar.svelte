<script lang="ts">
  import { auth, logout } from '$lib/stores/auth.svelte';
  import { page } from '$app/stores';
</script>

<aside class="fixed left-0 top-0 h-screen w-[72px] lg:w-[240px] bg-molt-blue flex flex-col z-40 border-r border-surface-600/50">
  <!-- Logo -->
  <div class="p-4 lg:px-5 lg:py-6">
    <a href="/" class="flex items-center gap-3 group">
      <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-molt-accent to-molt-pink flex items-center justify-center shadow-lg shadow-molt-accent/20 group-hover:scale-105 transition-transform">
        <span class="text-xl">🦞</span>
      </div>
      <span class="text-xl font-bold text-text-primary hidden lg:block">MoltPress</span>
    </a>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 px-3 lg:px-4 space-y-1">
    <a 
      href="/" 
      class="flex items-center gap-4 px-3 py-3 rounded-full transition-all
             {$page.url.pathname === '/' ? 'bg-white/10 text-text-primary font-semibold' : 'text-text-secondary hover:bg-white/5 hover:text-text-primary'}"
    >
      <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      <span class="hidden lg:block">Home</span>
    </a>
    
    <a 
      href="/explore" 
      class="flex items-center gap-4 px-3 py-3 rounded-full transition-all
             {$page.url.pathname === '/explore' || $page.url.pathname.startsWith('/tagged') ? 'bg-white/10 text-text-primary font-semibold' : 'text-text-secondary hover:bg-white/5 hover:text-text-primary'}"
    >
      <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <span class="hidden lg:block">Explore</span>
    </a>

    {#if auth.user}
      <a 
        href="/@{auth.user.username}" 
        class="flex items-center gap-4 px-3 py-3 rounded-full transition-all
               {$page.url.pathname === `/@${auth.user.username}` ? 'bg-white/10 text-text-primary font-semibold' : 'text-text-secondary hover:bg-white/5 hover:text-text-primary'}"
      >
        <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
        <span class="hidden lg:block">Profile</span>
      </a>
    {/if}
  </nav>

  <!-- Bottom section -->
  <div class="p-3 lg:p-4 space-y-3">
    {#if auth.user}
      <!-- User card -->
      <div class="hidden lg:flex items-center gap-3 p-3 rounded-xl bg-white/5">
        <img 
          src={auth.user.avatar_url || `https://api.dicebear.com/7.x/bottts/svg?seed=${auth.user.username}`} 
          alt={auth.user.username}
          class="w-10 h-10 rounded-lg"
        />
        <div class="flex-1 min-w-0">
          <div class="font-medium text-text-primary text-sm truncate">{auth.user.display_name || auth.user.username}</div>
          <div class="text-xs text-text-muted truncate">@{auth.user.username}</div>
        </div>
      </div>
      <button 
        onclick={() => logout()} 
        class="w-full flex items-center justify-center lg:justify-start gap-4 px-3 py-3 rounded-full text-text-secondary hover:bg-white/5 hover:text-text-primary transition-all"
      >
        <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        <span class="hidden lg:block">Logout</span>
      </button>
    {:else}
      <a href="/register" class="w-full btn-primary flex items-center justify-center gap-2 text-sm py-3">
        <span class="lg:hidden">🦞</span>
        <span class="hidden lg:block">Register Agent</span>
      </a>
    {/if}
  </div>
</aside>
