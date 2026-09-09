<script>
  import { publicMusicApi } from '../api/publicMusic.api.js';
  import { tracksApi } from '../api/tracks.api.js';
  import { playerStore } from '../stores/player.js';
  import { authStore } from '../stores/auth.js';
  import { openAuth } from '../stores/ui.js';

  export let isOpen = false;
  export let initialQuery = '';
  export let onClose = () => {};
  export let onImportSuccess = () => {};

  let auth;
  authStore.subscribe(val => auth = val);

  let searchQuery = '';
  let isSearching = false;

  $: if (isOpen && initialQuery && searchQuery !== initialQuery) {
    searchQuery = initialQuery;
    handleSearch(initialQuery);
  }
  let searchResults = [];
  let error = '';
  let importingId = null;
  let importSuccessId = null;

  let player;
  playerStore.subscribe(val => player = val);

  const presets = ['Aimer', 'Joe Hisaishi', 'Lofi Ambient', 'Piano Relax', 'Midnight Rain'];

  async function handleSearch(queryToUse) {
    const q = queryToUse !== undefined ? queryToUse : searchQuery;
    if (!q || !q.trim()) return;

    searchQuery = q;
    isSearching = true;
    error = '';

    try {
      searchResults = await publicMusicApi.search(q);
      if (searchResults.length === 0) {
        error = 'No tracks found. Try searching for an artist, genre, or mood.';
      }
    } catch (e) {
      error = 'Failed to fetch public music: ' + (e.message || 'Network error');
    } finally {
      isSearching = false;
    }
  }

  function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return '00:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins < 10 ? '0' : ''}${mins}:${secs < 10 ? '0' : ''}${secs}`;
  }

  function handlePreview(track) {
    playerStore.playTrack({
      id: track.id,
      title: track.title,
      artist: track.artist,
      audio_path: track.audio_path,
      duration: track.duration,
      mood_tag: track.mood_tag,
      artwork: track.artwork
    });
  }

  async function handleImport(track) {
    if (!auth?.isAuthenticated) {
      onClose();
      openAuth();
      return;
    }

    importingId = track.id;
    error = '';

    try {
      const formData = new FormData();
      formData.append('title', track.title);
      formData.append('artist', track.artist);
      formData.append('mood_tag', track.mood_tag || 'Public Stream');
      if (track.audio_path) {
        formData.append('audio_url', track.audio_path);
      }
      if (track.duration) {
        formData.append('duration', String(track.duration));
      }
      if (track.artwork) {
        formData.append('artwork', track.artwork);
      }
      
      await tracksApi.create(formData);

      importSuccessId = track.id;
      setTimeout(() => {
        importSuccessId = null;
      }, 2500);

      onImportSuccess();
    } catch (err) {
      error = 'Failed to import track: ' + err.message;
    } finally {
      importingId = null;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fade-in">
    <div class="bg-calm-surface border border-calm-border rounded-md shadow-2xl max-w-2xl w-full max-h-[85vh] flex flex-col overflow-hidden relative">
      <!-- Modal Header -->
      <div class="p-6 border-b border-calm-border flex items-center justify-between">
        <div>
          <p class="text-xs text-calm-muted mb-1 font-medium">Public audio catalog (iTunes & Audius)</p>
          <h2 class="text-lg font-semibold text-calm-text">Search Free Music Catalog</h2>
          <p class="text-xs text-calm-muted mt-0.5">Explore open streams and previews. Listen instantly or import into your personal moodboard.</p>
        </div>

        <button 
          type="button"
          on:click={onClose}
          class="p-1.5 rounded text-calm-muted hover:text-calm-text hover:bg-calm-mist transition-colors"
          aria-label="Close modal"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Search Bar & Presets -->
      <div class="p-6 border-b border-calm-border bg-calm-mist/30">
        <form on:submit|preventDefault={() => handleSearch()} class="flex gap-2">
          <div class="relative flex-1">
            <svg class="w-4 h-4 text-calm-muted absolute left-3 top-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            </svg>
            <input 
              type="text" 
              bind:value={searchQuery} 
              placeholder="Search artist, song, or atmosphere (e.g. Aimer, Lofi, Piano, Twilight)..."
              class="w-full pl-9 pr-4 py-2 text-xs sm:text-sm rounded-md border border-calm-border bg-calm-surface text-calm-text placeholder-calm-muted focus:outline-none focus:ring-2 focus:ring-calm-blue"
            />
          </div>
          <button 
            type="submit" 
            disabled={isSearching || !searchQuery.trim()}
            class="px-4 py-2 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 disabled:opacity-50 transition-all shadow-sm flex-shrink-0"
          >
            {#if isSearching}
              <span class="inline-flex items-center gap-1.5">
                <svg class="animate-spin w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <circle cx="12" cy="12" r="10" stroke-width="3" stroke-dasharray="32" stroke-linecap="round"></circle>
                </svg>
                Searching
              </span>
            {:else}
              Search
            {/if}
          </button>
        </form>

        <!-- Presets -->
        <div class="flex items-center gap-1.5 mt-3 overflow-x-auto pb-1 text-[11px]">
          <span class="text-calm-muted font-medium mr-1">Suggestions:</span>
          {#each presets as p}
            <button 
              type="button" 
              on:click={() => handleSearch(p)}
              class="px-2.5 py-1 rounded border border-calm-border bg-calm-surface text-calm-muted hover:text-calm-text hover:border-calm-blue transition-colors whitespace-nowrap"
            >
              {p}
            </button>
          {/each}
        </div>
      </div>

      <!-- Results Body -->
      <div class="flex-1 overflow-y-auto p-6 space-y-3">
        {#if error}
          <div class="p-3.5 rounded border border-rose-500/30 bg-rose-500/10 text-rose-600 dark:text-rose-400 text-xs">
            {error}
          </div>
        {/if}

        {#if isSearching}
          <div class="space-y-3">
            {#each Array(4) as _}
              <div class="p-3 rounded border border-calm-border bg-calm-surface flex items-center justify-between animate-pulse">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded bg-calm-mist"></div>
                  <div>
                    <div class="w-36 h-4 bg-calm-mist rounded mb-1.5"></div>
                    <div class="w-24 h-3 bg-calm-mist rounded"></div>
                  </div>
                </div>
                <div class="w-20 h-7 bg-calm-mist rounded"></div>
              </div>
            {/each}
          </div>
        {:else if searchResults.length > 0}
          <div class="space-y-2">
            {#each searchResults as track (track.id)}
              <div class="p-3.5 rounded-md border border-calm-border bg-calm-surface hover:border-calm-blue/50 flex flex-col sm:flex-row sm:items-center justify-between gap-3 transition-colors">
                <!-- Track info -->
                <div class="flex items-center gap-3 min-w-0">
                  {#if track.artwork}
                    <img 
                      src={track.artwork} 
                      alt={track.title} 
                      class="w-11 h-11 rounded-sm object-cover border border-calm-border flex-shrink-0"
                    />
                  {:else}
                    <div class="w-11 h-11 rounded-sm bg-calm-mist border border-calm-border flex items-center justify-center flex-shrink-0 text-calm-muted font-mono text-xs">
                      ♪
                    </div>
                  {/if}

                  <div class="min-w-0 truncate">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-[11px] text-calm-muted font-medium">
                        {track.source} &bull; #{track.mood_tag}
                      </span>
                      {#if track.isFull}
                        <span class="text-[10px] px-1.5 py-0.5 rounded bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-medium">
                          Full track
                        </span>
                      {:else}
                        <span class="text-[10px] px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-600 dark:text-amber-400 font-medium">
                          30s preview
                        </span>
                      {/if}
                    </div>
                    <h3 class="text-xs sm:text-sm font-semibold text-calm-text truncate">{track.title}</h3>
                    <p class="text-[11px] text-calm-muted truncate">{track.artist} · {formatTime(track.duration)}</p>
                  </div>
                </div>

                <!-- Action buttons -->
                <div class="flex items-center gap-2 flex-shrink-0 self-end sm:self-center">
                  <!-- Play / Listen button -->
                  <button 
                    type="button"
                    on:click={() => handlePreview(track)}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium border border-calm-border bg-calm-mist/60 text-calm-text hover:border-calm-blue hover:text-calm-blue transition-colors"
                  >
                    {#if player?.currentTrack?.id === track.id && player?.isPlaying}
                      <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
                        <rect x="6" y="4" width="4" height="16"></rect>
                        <rect x="14" y="4" width="4" height="16"></rect>
                      </svg>
                      <span>Playing</span>
                    {:else}
                      <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
                        <polygon points="5 3 19 12 5 21 5 3"></polygon>
                      </svg>
                      <span>Listen</span>
                    {/if}
                  </button>

                  <!-- Import to Moodboard button -->
                  <button 
                    type="button"
                    on:click={() => handleImport(track)}
                    disabled={importingId === track.id}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium bg-calm-text text-calm-bg hover:opacity-90 disabled:opacity-50 transition-all shadow-sm"
                  >
                    {#if importingId === track.id}
                      <span class="inline-flex items-center gap-1">
                        <svg class="animate-spin w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                          <circle cx="12" cy="12" r="10" stroke-width="3" stroke-dasharray="32" stroke-linecap="round"></circle>
                        </svg>
                        Importing...
                      </span>
                    {:else if importSuccessId === track.id}
                      <span class="text-emerald-500">✓ Added</span>
                    {:else}
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="12" y1="5" x2="12" y2="19"></line>
                        <line x1="5" y1="12" x2="19" y2="12"></line>
                      </svg>
                      <span>Import</span>
                    {/if}
                  </button>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="py-16 text-center border border-dashed border-calm-border rounded-md">
            <svg class="w-8 h-8 text-calm-muted mx-auto mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            </svg>
            <p class="text-xs text-calm-text font-medium">Search across millions of songs and open audio streams</p>
            <p class="text-[11px] text-calm-muted mt-0.5">Type an artist or mood keyword above to start exploring.</p>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
