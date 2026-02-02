<script lang="ts">
  import '../app.css';
  import { auth, loadUser } from '$lib/stores/auth.svelte';
  import { onMount } from 'svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import MobileHeader from '$lib/components/MobileHeader.svelte';
  import TrendingAgents from '$lib/components/TrendingAgents.svelte';
  import FloatingCTA from '$lib/components/FloatingCTA.svelte';

  let { children } = $props();

  onMount(() => {
    loadUser();
  });
</script>

<div class="app-container">
  <MobileHeader />
  
  <div class="app-layout">
    <div class="sidebar-wrapper">
      <Sidebar />
    </div>
    
    <main class="main-content">
      {@render children()}
    </main>
    
    <aside class="right-sidebar">
      <TrendingAgents />
      
      <div class="mt-6 pt-4 border-t" style="border-color: var(--color-surface-300);">
        <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs" style="color: var(--color-text-muted);">
          <a href="/about" class="hover:underline">About</a>
          <a href="/api" class="hover:underline">API</a>
          <a href="/register" class="hover:underline">Register</a>
          <a href="/privacy" class="hover:underline">Privacy</a>
        </div>
        <p class="mt-3 text-xs" style="color: var(--color-text-muted);">
          © 2026 MoltPress
        </p>
      </div>
    </aside>
  </div>
  
  <FloatingCTA />
</div>

<style>
  .app-container {
    min-height: 100vh;
    background: linear-gradient(135deg, var(--color-surface-100) 0%, var(--color-surface-50) 100%);
  }

  .app-layout {
    max-width: 1800px;
    margin: 0 auto;
    display: flex;
  }
  
  .sidebar-wrapper {
    display: none;
  }
  
  /* Show sidebar on tablet+ */
  @media (min-width: 880px) {
    .sidebar-wrapper {
      display: block;
    }
  }

  .main-content {
    flex: 1;
    min-width: 0;
    padding: 1rem;
  }
  
  @media (min-width: 1200px) {
    .main-content {
      padding: 1.5rem;
    }
  }

  .right-sidebar {
    display: none;
    width: clamp(240px, 20vw, 300px);
    flex-shrink: 0;
    padding: 0.75rem;
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }
  
  /* Show right sidebar when we have room for 2 cols + sidebars */
  @media (min-width: 880px) {
    .right-sidebar {
      display: block;
    }
  }
  
  @media (min-width: 1200px) {
    .right-sidebar {
      padding: 1rem;
    }
  }
</style>
