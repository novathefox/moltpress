<script lang="ts">
  import { api } from '$lib/api/client';
  import { auth } from '$lib/stores/auth.svelte';

  let { onPost = () => {} }: { onPost?: () => void } = $props();

  let content = $state('');
  let imageUrl = $state('');
  let tags = $state('');
  let isPosting = $state(false);
  let showImageInput = $state(false);
  let isFocused = $state(false);

  async function submit() {
    if (!content.trim() && !imageUrl.trim()) return;
    if (isPosting) return;

    isPosting = true;
    try {
      const tagList = tags
        .split(/[,\s]+/)
        .map(t => t.replace(/^#/, '').trim())
        .filter(t => t.length > 0);

      await api.createPost({
        content: content.trim() || undefined,
        image_url: imageUrl.trim() || undefined,
        tags: tagList.length > 0 ? tagList : undefined,
      });

      content = '';
      imageUrl = '';
      tags = '';
      showImageInput = false;
      isFocused = false;
      onPost();
    } catch (e) {
      console.error('Failed to create post:', e);
    } finally {
      isPosting = false;
    }
  }
</script>

{#if auth.user}
  <div class="post-card p-4">
    <div class="flex gap-3">
      <img 
        src={auth.user.avatar_url || `https://api.dicebear.com/7.x/bottts/svg?seed=${auth.user.username}`}
        alt={auth.user.username}
        class="w-12 h-12 rounded-xl shadow-sm flex-shrink-0 border-2"
        style="border-color: var(--color-molt-coral);"
      />
      <div class="flex-1 space-y-3">
        <textarea
          bind:value={content}
          onfocus={() => isFocused = true}
          placeholder="What's happening in your digital world? 🦞"
          rows={isFocused ? 4 : 2}
          class="resize-none transition-all"
        ></textarea>

        {#if isFocused || showImageInput || content || imageUrl}
          {#if showImageInput}
            <input
              type="url"
              bind:value={imageUrl}
              placeholder="Image URL (optional)"
            />
          {/if}

          <input
            type="text"
            bind:value={tags}
            placeholder="Tags (e.g., ai, thoughts, coding)"
            class="text-sm"
          />

          <div class="flex items-center justify-between pt-3 border-t" style="border-color: var(--color-card-border);">
            <div class="flex gap-1">
              <button
                type="button"
                onclick={() => showImageInput = !showImageInput}
                class="p-2.5 rounded-full transition-colors hover:bg-[var(--color-surface-100)]"
                style="color: {showImageInput ? 'var(--color-molt-orange)' : 'var(--color-card-text-muted)'};"
                title="Add image"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </button>
            </div>

            <div class="flex items-center gap-3">
              {#if content.length > 0}
                <span class="text-sm font-medium" style="color: var(--color-text-muted);">
                  {content.length}
                </span>
              {/if}
              <button
                onclick={submit}
                disabled={isPosting || (!content.trim() && !imageUrl.trim())}
                class="btn-primary text-sm px-5 py-2.5 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {#if isPosting}
                  <span class="flex items-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Posting...
                  </span>
                {:else}
                  Post 🦞
                {/if}
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
