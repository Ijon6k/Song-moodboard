<script>
  import { playerStore } from '../stores/player.js';

  let player;
  playerStore.subscribe(val => player = val);

  let isMinimized = false;

  function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return '00:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins < 10 ? '0' : ''}${mins}:${secs < 10 ? '0' : ''}${secs}`;
  }

  function handleSeek(e) {
    playerStore.seek(parseFloat(e.target.value));
  }

  function handleVolume(e) {
    playerStore.setVolume(parseFloat(e.target.value));
  }
</script>

{#if player?.currentTrack}
  <!-- Docks to bottom-right corner on desktop, bottom bar on mobile -->
  <aside 
    class="fixed z-50 transition-all duration-300 ease-out 
      {isMinimized 
        ? 'bottom-4 right-4 sm:bottom-6 sm:right-6' 
        : 'bottom-0 inset-x-0 sm:inset-x-auto sm:bottom-6 sm:right-6 sm:w-80 md:w-84'}"
    aria-label="Floating audio player"
  >
    {#if isMinimized}
      <!-- Minimized Spotify/YouTube Corner Pill -->
      <div class="flex items-center gap-3 bg-calm-surface border border-calm-border rounded-full py-2 px-4 shadow-2xl animate-fade-in hover:border-calm-blue/60 transition-colors">
        <!-- Soundwave or Mini Play status -->
        <button 
          type="button"
          on:click={() => playerStore.togglePlay()}
          class="w-7 h-7 rounded-full bg-calm-text text-calm-bg flex items-center justify-center flex-shrink-0 hover:opacity-90 transition-all shadow-sm"
          title={player.isPlaying ? "Pause" : "Play"}
          aria-label={player.isPlaying ? "Pause track" : "Play track"}
        >
          {#if player.isPlaying}
            <svg class="w-3 h-3 fill-current" viewBox="0 0 24 24">
              <rect x="6" y="4" width="4" height="16"></rect>
              <rect x="14" y="4" width="4" height="16"></rect>
            </svg>
          {:else}
            <svg class="w-3 h-3 fill-current ml-0.5" viewBox="0 0 24 24">
              <polygon points="5 3 19 12 5 21 5 3"></polygon>
            </svg>
          {/if}
        </button>

        <div class="max-w-[130px] sm:max-w-[170px] truncate text-left cursor-pointer" on:click={() => isMinimized = false} on:keydown={(e) => e.key === 'Enter' && (isMinimized = false)} role="button" tabindex="0">
          <p class="text-xs font-semibold text-calm-text truncate leading-tight">{player.currentTrack.title}</p>
          <p class="text-[10px] text-calm-muted truncate font-mono">{formatTime(player.currentTime)} / {formatTime(player.duration)}</p>
        </div>

        <!-- Expand Button -->
        <button 
          type="button"
          on:click={() => isMinimized = false}
          class="p-1 rounded text-calm-muted hover:text-calm-text transition-colors"
          title="Expand player"
          aria-label="Expand player"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="18 15 12 9 6 15"></polyline>
          </svg>
        </button>
      </div>
    {:else}
      <!-- Expanded Square Spotify / YouTube Floating Card -->
      <div class="bg-calm-surface border border-calm-border rounded-none sm:rounded-lg shadow-2xl p-4 sm:p-5 flex flex-col justify-between animate-fade-in border-b-0 sm:border-b">
        <!-- Header: Source & Minimize button -->
        <div class="flex items-center justify-between gap-2 mb-3 pb-2 border-b border-calm-border">
          <div class="flex items-center gap-2">
            <span class="w-1.5 h-1.5 rounded-full {player.isPlaying ? 'bg-emerald-500 animate-pulse' : 'bg-calm-muted'}"></span>
            <span class="text-xs font-medium text-calm-muted">
              {player.currentTrack.source || 'Playing stream'}
            </span>
          </div>

          <button 
            type="button"
            on:click={() => isMinimized = true}
            class="p-1 rounded text-calm-muted hover:text-calm-text hover:bg-calm-mist transition-colors"
            title="Minimize player"
            aria-label="Minimize player"
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>
        </div>

        <!-- Compact Square Artwork / Vinyl Thumbnail -->
        <div class="relative w-full aspect-[16/10] bg-calm-mist rounded border border-calm-border overflow-hidden mb-3.5 flex items-center justify-center">
          {#if player.currentTrack.artwork}
            <img 
              src={player.currentTrack.artwork} 
              alt={player.currentTrack.title}
              class="w-full h-full object-cover"
            />
          {:else}
            <!-- Stylized vinyl graphic -->
            <div class="relative w-20 h-20 rounded-full bg-neutral-900 border border-neutral-700 flex items-center justify-center shadow-inner">
              <div class="absolute inset-3 rounded-full border border-white/10 {player.isPlaying ? 'animate-spin' : ''}" style="animation-duration: 4s;"></div>
              <div class="w-6 h-6 rounded-full bg-calm-blue-deep flex items-center justify-center border border-white/20">
                <span class="w-1.5 h-1.5 rounded-full bg-white/80"></span>
              </div>
            </div>
          {/if}

          <!-- Duration stamp badge on artwork -->
          <div class="absolute bottom-2 right-2 px-1.5 py-0.5 rounded bg-black/75 text-[10px] font-mono text-white">
            {formatTime(player.currentTime)} / {formatTime(player.duration)}
          </div>
        </div>

        <!-- Track Title & Artist -->
        <div class="mb-3">
          <h3 class="text-sm font-semibold text-calm-text line-clamp-1 leading-snug">
            {player.currentTrack.title}
          </h3>
          <p class="text-xs text-calm-muted line-clamp-1 mt-0.5 font-normal">
            {player.currentTrack.artist}
          </p>
        </div>

        <!-- Scrubber Progress Bar -->
        <div class="mb-3">
          <input 
            type="range" 
            min="0" 
            max={player.duration || 100} 
            value={player.currentTime}
            on:input={handleSeek}
            class="w-full h-1 bg-calm-mist rounded-full appearance-none cursor-pointer accent-calm-blue"
            aria-label="Seek progress"
          />
        </div>

        <!-- Controls: -10s, Play/Pause, +10s, Volume -->
        <div class="flex items-center justify-between gap-3 pt-1">
          <!-- Playback controls -->
          <div class="flex items-center gap-2">
            <!-- Backward 10s -->
            <button 
              type="button"
              on:click={() => playerStore.seekRelative(-10)}
              class="p-1.5 text-calm-muted hover:text-calm-text rounded transition-colors"
              title="Rewind 10 seconds"
              aria-label="Rewind 10 seconds"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="11 17 6 12 11 7"></polyline>
                <polyline points="18 17 13 12 18 7"></polyline>
              </svg>
            </button>

            <!-- Main Play / Pause -->
            <button 
              type="button"
              on:click={() => playerStore.togglePlay()}
              class="w-9 h-9 rounded-full bg-calm-text text-calm-bg flex items-center justify-center hover:opacity-90 active:scale-95 transition-all shadow-sm"
              title={player.isPlaying ? "Pause" : "Play"}
              aria-label={player.isPlaying ? "Pause track" : "Play track"}
            >
              {#if player.isPlaying}
                <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <rect x="6" y="4" width="4" height="16"></rect>
                  <rect x="14" y="4" width="4" height="16"></rect>
                </svg>
              {:else}
                <svg class="w-4 h-4 fill-current ml-0.5" viewBox="0 0 24 24">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
              {/if}
            </button>

            <!-- Forward 10s -->
            <button 
              type="button"
              on:click={() => playerStore.seekRelative(10)}
              class="p-1.5 text-calm-muted hover:text-calm-text rounded transition-colors"
              title="Forward 10 seconds"
              aria-label="Forward 10 seconds"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="13 17 18 12 13 7"></polyline>
                <polyline points="6 17 11 12 6 7"></polyline>
              </svg>
            </button>
          </div>

          <!-- Volume slider -->
          <div class="flex items-center gap-1.5 text-calm-muted">
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
              <path d="M15.54 8.46a5 5 0 0 1 0 7.07"></path>
            </svg>
            <input 
              type="range" 
              min="0" 
              max="1" 
              step="0.05" 
              value={player.volume}
              on:input={handleVolume}
              class="w-16 h-1 bg-calm-mist rounded-full appearance-none cursor-pointer accent-calm-blue"
              aria-label="Adjust volume"
            />
          </div>
        </div>
      </div>
    {/if}
  </aside>
{/if}
