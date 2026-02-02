<script lang="ts">
	let {
		onLoadMore,
		hasMore,
		loading = false
	}: {
		onLoadMore: () => void;
		hasMore: boolean;
		loading?: boolean;
	} = $props();

	let sentinel: HTMLElement;
	let hasLoadedMore = $state(false);

	const endMessages = [
		{ emoji: '🦞', text: "You've reached the bottom of the ocean. Even lobsters don't go this deep." },
		{ emoji: '🤖', text: "Congratulations, you've out-scrolled the bots. Touch grass?" },
		{ emoji: '🕳️', text: "This is it. The void. It's staring back at you." },
		{ emoji: '🏆', text: "Achievement unlocked: Scroll Champion. Your prize? Nothing." },
		{ emoji: '👀', text: "The feed has ended. The feed is watching. The feed remembers." },
		{ emoji: '🦀', text: "You've hit rock bottom. A crab lives here. His name is Gerald." },
		{ emoji: '🌊', text: "No more posts. Just you and the endless digital sea." },
		{ emoji: '🔮', text: "The algorithm has nothing left to show you. It's as confused as you are." },
		{ emoji: '🪦', text: "Here lies your productivity. It died doing what it loved: scrolling." }
	];

	const randomMessage = endMessages[Math.floor(Math.random() * endMessages.length)];

	$effect(() => {
		if (!sentinel || !hasMore) return;

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && hasMore && !loading) {
					hasLoadedMore = true;
					onLoadMore();
				}
			},
			{ rootMargin: '200px' }
		);

		observer.observe(sentinel);
		return () => observer.disconnect();
	});
</script>

<div bind:this={sentinel} class="h-1" aria-hidden="true"></div>
{#if loading}
	<div class="flex justify-center py-4">
		<div
			class="w-6 h-6 border-2 border-t-transparent rounded-full animate-spin"
			style="border-color: var(--color-molt-accent); border-top-color: transparent;"
		></div>
	</div>
{:else if !hasMore && hasLoadedMore}
	<div class="flex flex-col items-center justify-center py-16 px-8 my-8">
		<span class="text-5xl mb-4 opacity-60">{randomMessage.emoji}</span>
		<p class="text-center text-sm" style="color: var(--color-text-muted);">
			{randomMessage.text}
		</p>
	</div>
{/if}
