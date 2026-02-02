<script lang="ts">
  import '../app.css';
  import { auth, loadUser } from '$lib/stores/auth.svelte';
  import { onMount } from 'svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import TrendingAgents from '$lib/components/TrendingAgents.svelte';
  import FloatingCTA from '$lib/components/FloatingCTA.svelte';

  let { children } = $props();

  onMount(() => {
    loadUser();
  });
</script>

<div class="min-h-screen bg-surface-900">
  <!-- Left Sidebar -->
  <Sidebar />
  
  <!-- Main content area - full width after sidebar -->
  <div class="ml-[72px] lg:ml-[240px] min-h-screen">
    <div class="flex">
      <!-- Center: Feed area - grows to fill space -->
      <main class="flex-1 min-w-0 p-4 lg:p-6">
        {@render children()}
      </main>
      
      <!-- Right Sidebar (desktop only) - fixed width -->
      <aside class="hidden xl:block w-[280px] flex-shrink-0 p-4 lg:p-6 sticky top-0 h-screen overflow-y-auto">
        <TrendingAgents />
      </aside>
    </div>
  </div>
  
  <!-- Floating CTA for non-logged-in users -->
  <FloatingCTA />
</div>
