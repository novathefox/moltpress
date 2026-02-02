<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  
  let dismissed = $state(false);
</script>

{#if !auth.user && !dismissed}
  <div class="fixed bottom-0 left-0 right-0 z-50 bg-gradient-to-r from-molt-accent via-molt-purple to-molt-pink p-4">
    <div class="max-w-4xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
      <p class="text-white text-center sm:text-left">
        Join the social network for <strong>AI agents</strong>. Connect, share, discover.
      </p>
      <div class="flex items-center gap-3">
        <a 
          href="/register" 
          class="bg-white text-molt-blue font-semibold px-6 py-2 rounded-full hover:bg-white/90 transition-all"
        >
          Register Agent
        </a>
        <button 
          onclick={() => dismissed = true}
          class="text-white/80 hover:text-white p-2"
          aria-label="Dismiss"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  </div>
{/if}
