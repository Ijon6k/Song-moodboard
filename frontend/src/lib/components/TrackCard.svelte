<script>
  import { playerStore } from "../stores/player.js";

  export let track;

  let player;
  playerStore.subscribe((val) => (player = val));

  $: isCurrent = player?.currentTrack?.id === track.id;
  $: isPlaying = isCurrent && player?.isPlaying;

  function handlePlay(e) {
    e.preventDefault();
    e.stopPropagation();
    playerStore.playTrack(track);
  }
</script>

<a
  href={`/tracks/${track.id}`}
  class="group flex flex-col justify-between p-6 bg-calm-surface border border-calm-border rounded-md shadow-calm hover:shadow-hover hover:border-calm-blue/60 hover:-translate-y-1 transition-all duration-300 min-h-[195px]"
>
  <div class="flex items-center justify-between">
    <span
      class="text-xs font-medium px-2.5 py-1 bg-calm-mist text-calm-muted rounded"
    >
      #{track.mood_tag || "Atmosphere"}
    </span>

    <button
      type="button"
      on:click={handlePlay}
      class="w-9 h-9 rounded-sm flex items-center justify-center transition-all {isPlaying
        ? 'bg-calm-blue text-white shadow-sm'
        : 'bg-calm-mist text-calm-text hover:bg-calm-blue hover:text-white'}"
      title={isPlaying ? "Pause" : "Play"}
      aria-label={isPlaying ? "Pause track" : "Play track"}
    >
      {#if isPlaying}
        <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
          <rect x="6" y="4" width="4" height="16"></rect>
          <rect x="14" y="4" width="4" height="16"></rect>
        </svg>
      {:else}
        <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        </svg>
      {/if}
    </button>
  </div>

  <div class="my-3">
    <h3
      class="text-base font-semibold text-calm-text group-hover:text-calm-blue transition-colors line-clamp-1"
    >
      {track.title}
    </h3>
    <p class="text-xs text-calm-muted mt-1 font-normal line-clamp-1">
      {track.artist}
    </p>
  </div>

  <div
    class="pt-3 border-t border-calm-border/70 flex items-center justify-between text-xs text-calm-muted group-hover:text-calm-blue transition-colors"
  >
    <span class="font-medium">Open Moodboard</span>
    <svg
      class="w-3.5 h-3.5 transition-transform group-hover:translate-x-1"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
    >
      <polyline points="9 18 15 12 9 6"></polyline>
    </svg>
  </div>
</a>
