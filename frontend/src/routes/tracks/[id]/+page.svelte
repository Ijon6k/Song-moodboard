<script>
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount, onDestroy } from 'svelte';
  import { createQuery, useQueryClient } from '@tanstack/svelte-query';
  import { tracksApi } from '$lib/api/tracks.api.js';
  import { queryKeys } from '$lib/queries/keys.js';
  import { playerStore } from '$lib/stores/player.js';
  import { authStore } from '$lib/stores/auth.js';
  import { openAuth } from '$lib/stores/ui.js';
  import AddMoodItemModal from '$lib/components/AddMoodItemModal.svelte';

  const queryClient = useQueryClient();
  $: trackId = $page.params.id;

  let auth;
  authStore.subscribe(val => auth = val);

  $: trackQuery = createQuery({
    queryKey: queryKeys.tracks.detail(trackId),
    queryFn: () => tracksApi.getById(trackId),
    enabled: !!trackId && !!auth?.isAuthenticated,
  });

  $: track = $trackQuery.data;

  async function handleDeleteTrack() {
    if (!track) return;
    if (!confirm(`Hapus lagu "${track.title}" dari perpustakaan?`)) return;
    try {
      if (isCurrent) {
        playerStore.pause();
      }
      await tracksApi.delete(track.id);
      queryClient.invalidateQueries({ queryKey: queryKeys.tracks.all });
      goto('/');
    } catch (err) {
      alert('Gagal menghapus lagu: ' + (err.message || 'Terjadi kesalahan'));
    }
  }

  let player;
  playerStore.subscribe(val => player = val);

  $: isCurrent = player?.currentTrack?.id === track?.id;
  $: isPlaying = isCurrent && player?.isPlaying;

  let isAddMoodOpen = false;
  let activeTab = 'all';

  // Lightbox Zoom state
  let activeZoomItem = null;
  let zoomScale = 1;

  function openZoom(item) {
    activeZoomItem = item;
    zoomScale = 1;
  }

  function closeZoom() {
    activeZoomItem = null;
    zoomScale = 1;
  }

  function zoomIn() {
    zoomScale = Math.min(3, zoomScale + 0.3);
  }

  function zoomOut() {
    zoomScale = Math.max(0.7, zoomScale - 0.3);
  }

  function resetZoom() {
    zoomScale = 1;
  }

  function handleKeydown(e) {
    if (e.key === 'Escape' && activeZoomItem) {
      closeZoom();
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
  });

  onDestroy(() => {
    if (typeof window !== 'undefined') {
      window.removeEventListener('keydown', handleKeydown);
    }
  });

  function togglePlay() {
    if (track) {
      playerStore.playTrack(track);
    }
  }

  function handleSeek(e) {
    if (isCurrent) {
      playerStore.seek(parseFloat(e.target.value));
    } else {
      playerStore.playTrack(track);
    }
  }

  function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return '00:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins < 10 ? '0' : ''}${mins}:${secs < 10 ? '0' : ''}${secs}`;
  }

  function handleAddMoodSuccess() {
    queryClient.invalidateQueries({ queryKey: queryKeys.tracks.detail(trackId) });
  }

  $: items = track?.moodboard_items || [];
  $: counts = {
    all: items.length,
    image: items.filter(it => it.item_type === 'image').length,
    quote: items.filter(it => it.item_type === 'quote').length,
    palette: items.filter(it => it.item_type === 'palette').length,
    note: items.filter(it => it.item_type === 'note').length,
  };

  $: filteredItems = activeTab === 'all' 
    ? items 
    : items.filter(it => it.item_type === activeTab);
</script>

<div class="max-w-[1340px] mx-auto px-6 sm:px-10 py-10 animate-fade-in text-calm-text">
  <!-- Top Navigation & Breadcrumbs (Clean, no fake badges) -->
  <div class="flex items-center justify-between gap-4 mb-6">
    <a 
      href="/" 
      class="inline-flex items-center gap-2 text-xs text-calm-muted hover:text-calm-text transition-colors group"
    >
      <svg class="w-4 h-4 transition-transform group-hover:-translate-x-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
        <line x1="19" y1="12" x2="5" y2="12"></line>
        <polyline points="12 19 5 12 12 5"></polyline>
      </svg>
      <span>Back to tracks</span>
    </a>

    {#if track?.mood_tag}
      <span class="text-xs text-calm-muted">
        Category: <strong class="text-calm-text font-medium">{track.mood_tag}</strong>
      </span>
    {/if}
  </div>

  {#if !auth?.isAuthenticated}
    <div class="p-12 text-center bg-calm-surface border border-calm-border rounded-xl shadow-sm my-8">
      <div class="w-12 h-12 rounded-full bg-calm-ice text-calm-blue mx-auto flex items-center justify-center mb-3">
        <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
          <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
        </svg>
      </div>
      <h2 class="text-lg font-semibold text-calm-text mb-1">Login Diperlukan</h2>
      <p class="text-xs sm:text-sm text-calm-muted max-w-sm mx-auto mb-6">
        Lagu dan moodboard ini bersifat privat. Silakan login dengan akun pemilik untuk melihat atau mendengarkan.
      </p>
      <div class="flex items-center justify-center gap-3">
        <button
          type="button"
          on:click={openAuth}
          class="px-5 py-2.5 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 transition-all shadow-sm"
        >
          Masuk Akun
        </button>
        <a
          href="/"
          class="px-4 py-2.5 rounded-md border border-calm-border bg-calm-surface text-calm-text text-xs sm:text-sm font-medium hover:bg-calm-mist transition-colors shadow-sm"
        >
          Kembali ke Beranda
        </a>
      </div>
    </div>
  {:else if $trackQuery?.isLoading}
    <div class="h-44 bg-calm-surface border border-calm-border rounded-lg animate-pulse mb-8"></div>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {#each Array(8) as _}
        <div class="h-64 bg-calm-surface border border-calm-border rounded-lg animate-pulse"></div>
      {/each}
    </div>
  {:else if $trackQuery?.isError}
    <div class="p-12 text-center bg-calm-surface border border-rose-500/20 rounded-xl text-rose-600 dark:text-rose-400 text-sm shadow-sm my-8">
      <p class="font-semibold text-base mb-1">Tidak Dapat Mengakses Lagu</p>
      <p class="text-xs opacity-80 mb-5">{$trackQuery.error?.message || 'Lagu tidak ditemukan atau Anda tidak memiliki akses ke lagu ini.'}</p>
      <a
        href="/"
        class="inline-block px-4 py-2 rounded-md bg-calm-text text-calm-bg text-xs font-medium hover:opacity-90 transition-all shadow-sm"
      >
        Kembali ke Perpustakaan Saya
      </a>
    </div>
  {:else if track}
    <!-- Track Master Header: Minimalist, Balanced, No Capslock Slop -->
    <div class="bg-calm-surface border border-calm-border rounded-lg p-6 sm:p-8 shadow-sm mb-10">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <!-- Track Visual & Information -->
        <div class="flex items-center gap-5 min-w-0">
          <!-- Realistic Vinyl Disc with Subtle Grooves -->
          <div class="relative w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-neutral-900 border border-neutral-700 flex items-center justify-center flex-shrink-0 shadow-sm overflow-hidden">
            <div class="absolute inset-2 rounded-full border border-white/10 {isPlaying ? 'animate-spin' : ''}" style="animation-duration: 4s;"></div>
            <div class="absolute inset-4 rounded-full border border-white/10 {isPlaying ? 'animate-spin' : ''}" style="animation-duration: 4s;"></div>
            <div class="relative w-6 h-6 rounded-full bg-calm-blue flex items-center justify-center border border-white/20 z-10">
              <span class="w-1.5 h-1.5 rounded-full {isPlaying ? 'bg-white animate-ping' : 'bg-white/80'}"></span>
            </div>
          </div>

          <!-- Metadata -->
          <div class="min-w-0">
            <h1 class="text-xl sm:text-2xl md:text-3xl font-semibold tracking-tight text-calm-text truncate leading-tight">
              {track.title}
            </h1>
            <p class="text-sm text-calm-muted mt-1 truncate">
              {track.artist}
            </p>
          </div>
        </div>

        <!-- Controls & Actions -->
        <div class="flex items-center gap-3 flex-shrink-0">
          <!-- Play / Pause -->
          <button 
            type="button"
            on:click={togglePlay}
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-md text-sm font-medium transition-all shadow-sm {isPlaying ? 'bg-calm-blue text-white' : 'bg-calm-text text-calm-bg hover:opacity-90'}"
          >
            {#if isPlaying}
              <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <rect x="6" y="4" width="4" height="16"></rect>
                <rect x="14" y="4" width="4" height="16"></rect>
              </svg>
              <span>Pause</span>
            {:else}
              <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
              <span>Play track</span>
            {/if}
          </button>

          <!-- Add Fragment Button -->
          <button 
            type="button"
            on:click={() => isAddMoodOpen = true}
            class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-md border border-calm-border bg-calm-surface text-calm-text text-sm font-medium hover:bg-calm-mist transition-colors shadow-sm"
          >
            <svg class="w-4 h-4 text-calm-blue" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
            <span>Add fragment</span>
          </button>

          <!-- Delete Track Button -->
          <button 
            type="button"
            on:click={handleDeleteTrack}
            class="inline-flex items-center gap-1.5 px-3.5 py-2.5 rounded-md border border-calm-border bg-calm-surface text-neutral-400 hover:text-rose-600 hover:border-rose-500/30 hover:bg-rose-500/5 text-sm font-medium transition-colors shadow-sm"
            title="Delete this track"
            aria-label="Delete this track"
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            <span class="hidden sm:inline">Delete</span>
          </button>
        </div>
      </div>

      <!-- Scrubber Bar -->
      <div class="mt-6 pt-5 border-t border-calm-border flex items-center gap-3">
        <span class="text-xs text-calm-muted w-10 text-left">
          {isCurrent ? formatTime(player.currentTime) : '00:00'}
        </span>
        <div class="flex-1 relative flex items-center">
          <input 
            type="range" 
            min="0" 
            max={isCurrent && player.duration ? player.duration : (track.duration || 100)} 
            value={isCurrent ? player.currentTime : 0}
            on:input={handleSeek}
            class="w-full h-1 bg-calm-mist rounded-full appearance-none cursor-pointer accent-calm-blue"
            aria-label="Track progress"
          />
        </div>
        <span class="text-xs text-calm-muted w-10 text-right">
          {isCurrent && player.duration ? formatTime(player.duration) : formatTime(track.duration || 180)}
        </span>
      </div>
    </div>

    <!-- Moodboard Canvas Header & Category Tabs -->
    <section>
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6 pb-4 border-b border-calm-border">
        <div>
          <h2 class="text-lg font-semibold text-calm-text">Moodboard fragments</h2>
          <p class="text-xs text-calm-muted mt-0.5">Photographs, lyrics, and colors curated for this melody.</p>
        </div>

        <!-- Filter tabs with fragment counts -->
        <div class="flex items-center gap-1 bg-calm-mist/50 p-1 rounded-md border border-calm-border overflow-x-auto">
          {#each [
            { id: 'all', label: 'All', count: counts.all },
            { id: 'image', label: 'Photos', count: counts.image },
            { id: 'quote', label: 'Lyrics', count: counts.quote },
            { id: 'palette', label: 'Palettes', count: counts.palette },
            { id: 'note', label: 'Notes', count: counts.note }
          ] as tab}
            <button 
              type="button"
              on:click={() => activeTab = tab.id}
              class="px-3 py-1.5 text-xs font-medium rounded transition-all whitespace-nowrap {activeTab === tab.id ? 'bg-calm-surface text-calm-text shadow-sm' : 'text-calm-muted hover:text-calm-text'}"
            >
              <span>{tab.label}</span>
              <span class="text-[11px] opacity-60 ml-0.5">({tab.count})</span>
            </button>
          {/each}
        </div>
      </div>

      {#if filteredItems.length === 0}
        <div class="p-16 text-center bg-calm-surface border border-dashed border-calm-border rounded-lg">
          <p class="text-sm font-medium text-calm-text mb-1">No fragments in this tab</p>
          <p class="text-xs text-calm-muted mb-4">Attach a photo, quote, or note to build this moodboard.</p>
          <button 
            type="button"
            on:click={() => isAddMoodOpen = true}
            class="px-4 py-2 text-xs font-medium rounded-md bg-calm-text text-calm-bg hover:opacity-90 shadow-sm transition-all"
          >
            Add fragment
          </button>
        </div>
      {:else}
        <!-- Disciplined, Aligned Responsive CSS Grid (Fast thumb loading + Zoom trigger) -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 items-stretch">
          {#each filteredItems as item (item.id)}
            {#if item.item_type === 'image'}
              <!-- Photo Fragment: Uses lightweight thumbnail, click to zoom original -->
              <button 
                type="button"
                on:click={() => openZoom(item)}
                class="text-left bg-calm-surface border border-calm-border rounded-lg overflow-hidden flex flex-col justify-between hover:border-calm-blue/60 transition-all shadow-sm group focus:outline-none focus:ring-2 focus:ring-calm-blue"
                title="Click to view original image in full resolution"
              >
                <!-- Thumbnail Frame -->
                <div class="relative w-full aspect-[4/3] bg-calm-mist overflow-hidden">
                  <img 
                    src={item.thumb_path || item.media_path} 
                    alt={item.content || "Mood photo"} 
                    loading="lazy" 
                    class="w-full h-full object-cover group-hover:scale-[1.02] transition-transform duration-300"
                  />
                  <!-- Subtle Zoom Hint Badge -->
                  <div class="absolute bottom-2 right-2 p-1.5 rounded-md bg-black/60 text-white opacity-0 group-hover:opacity-100 transition-opacity">
                    <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="11" cy="11" r="8"></circle>
                      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                      <line x1="11" y1="8" x2="11" y2="14"></line>
                      <line x1="8" y1="11" x2="14" y2="11"></line>
                    </svg>
                  </div>
                </div>

                <div class="p-3.5 flex-1 flex flex-col justify-between border-t border-calm-border">
                  <p class="text-xs text-calm-muted leading-relaxed font-normal">
                    {item.content || "Photograph fragment"}
                  </p>
                  <span class="text-[11px] text-calm-blue font-medium mt-2 flex items-center gap-1">
                    <span>Click to zoom</span>
                    <span>→</span>
                  </span>
                </div>
              </button>

            {:else if item.item_type === 'quote'}
              <!-- Quote Fragment: Editorial typography, no left-stripe slop -->
              <div class="bg-calm-surface border border-calm-border rounded-lg p-5 flex flex-col justify-between hover:border-calm-blue/60 transition-colors shadow-sm">
                <div>
                  <span class="text-3xl leading-none text-calm-blue font-serif block mb-1 opacity-70">“</span>
                  <p class="text-sm font-serif italic text-calm-text leading-relaxed">
                    {item.content}
                  </p>
                </div>
                <div class="pt-3 mt-4 border-t border-calm-border flex items-center justify-between text-xs text-calm-muted">
                  <span>Lyrics</span>
                </div>
              </div>

            {:else if item.item_type === 'palette'}
              <!-- Palette Fragment: Clean color swatches -->
              <div class="bg-calm-surface border border-calm-border rounded-lg p-5 flex flex-col justify-between hover:border-calm-blue/60 transition-colors shadow-sm">
                <div>
                  <div class="flex items-center justify-between mb-3 text-xs font-medium text-calm-text">
                    <span>Palette</span>
                  </div>
                  <div class="grid grid-cols-5 gap-1.5">
                    {#each item.content.split(',') as color}
                      <div class="flex flex-col items-center gap-1.5">
                        <div class="w-full aspect-square rounded border border-black/10 shadow-sm" style="background-color: {color.trim()};"></div>
                        <span class="text-[9px] text-calm-muted truncate w-full text-center">{color.trim()}</span>
                      </div>
                    {/each}
                  </div>
                </div>
                <div class="pt-3 mt-4 border-t border-calm-border text-xs text-calm-muted">
                  Harmonized tones
                </div>
              </div>

            {:else}
              <!-- Reflection / Note Fragment: Clean memo styling -->
              <div class="bg-calm-surface border border-calm-border rounded-lg p-5 flex flex-col justify-between hover:border-calm-blue/60 transition-colors shadow-sm">
                <div>
                  <span class="text-xs font-medium text-calm-text block mb-2">
                    Note
                  </span>
                  <p class="text-xs text-calm-text/85 leading-relaxed font-normal">
                    {item.content}
                  </p>
                </div>
                <div class="pt-3 mt-4 border-t border-calm-border text-xs text-calm-muted">
                  Personal reflection
                </div>
              </div>
            {/if}
          {/each}
        </div>
      {/if}
    </section>
  {/if}

  <!-- Add Item Modal -->
  <AddMoodItemModal 
    isOpen={isAddMoodOpen} 
    {trackId}
    onClose={() => isAddMoodOpen = false}
    onSuccess={handleAddMoodSuccess}
  />

  <!-- Interactive Image Zoom Lightbox Modal (Original Full-Res Display) -->
  {#if activeZoomItem}
    <div 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-8 animate-fade-in"
      role="dialog"
      aria-modal="true"
      aria-label="Image Zoom Lightbox"
    >
      <!-- Backdrop button for closing -->
      <button 
        type="button" 
        class="absolute inset-0 bg-black/90 backdrop-blur-md cursor-default w-full h-full border-0 p-0 m-0" 
        on:click={closeZoom} 
        aria-label="Close lightbox"
      ></button>

      <!-- Floating Toolbar at top -->
      <div 
        class="absolute top-5 left-1/2 -translate-x-1/2 z-10 flex items-center gap-2 bg-neutral-900/90 border border-neutral-700 rounded-full px-4 py-1.5 text-white shadow-2xl"
      >
        <!-- Zoom out -->
        <button 
          type="button" 
          on:click={zoomOut}
          disabled={zoomScale <= 0.7}
          class="p-1.5 hover:text-sky-400 disabled:opacity-40 transition-colors"
          title="Zoom out (-)"
          aria-label="Zoom out"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            <line x1="8" y1="11" x2="14" y2="11"></line>
          </svg>
        </button>

        <span class="text-xs font-mono px-1 select-none">{Math.round(zoomScale * 100)}%</span>

        <!-- Zoom in -->
        <button 
          type="button" 
          on:click={zoomIn}
          disabled={zoomScale >= 3}
          class="p-1.5 hover:text-sky-400 disabled:opacity-40 transition-colors"
          title="Zoom in (+)"
          aria-label="Zoom in"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            <line x1="11" y1="8" x2="11" y2="14"></line>
            <line x1="8" y1="11" x2="14" y2="11"></line>
          </svg>
        </button>

        <span class="w-px h-4 bg-neutral-700 mx-1"></span>

        <!-- Reset 1:1 -->
        <button 
          type="button" 
          on:click={resetZoom}
          class="text-xs hover:text-sky-400 px-2 py-0.5 rounded transition-colors"
          title="Reset zoom"
        >
          Reset
        </button>

        <span class="w-px h-4 bg-neutral-700 mx-1"></span>

        <!-- Close -->
        <button 
          type="button" 
          on:click={closeZoom}
          class="p-1.5 hover:text-rose-400 transition-colors"
          title="Close (ESC)"
          aria-label="Close lightbox"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Image Canvas Container (Displays High-Res Original) -->
      <div 
        class="relative z-0 max-w-full max-h-full flex items-center justify-center overflow-auto p-4 pointer-events-none"
      >
        <img 
          src={activeZoomItem.media_path} 
          alt={activeZoomItem.content || "Original high-res photo"}
          style="transform: scale({zoomScale}); transition: transform 0.2s ease-out;"
          class="max-w-[90vw] max-h-[82vh] object-contain rounded shadow-2xl select-none pointer-events-auto"
        />
      </div>

      <!-- Caption at bottom -->
      {#if activeZoomItem.content}
        <div 
          class="absolute bottom-6 left-1/2 -translate-x-1/2 max-w-xl text-center px-4 py-2 rounded-md bg-neutral-900/80 text-white text-xs backdrop-blur-sm border border-neutral-700 z-10 pointer-events-none"
        >
          {activeZoomItem.content}
        </div>
      {/if}
    </div>
  {/if}
</div>
