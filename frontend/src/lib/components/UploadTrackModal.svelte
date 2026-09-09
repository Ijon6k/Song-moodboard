<script>
  import { onDestroy } from 'svelte';
  import { tracksApi } from '../api/tracks.api.js';
  import { categoryStore } from '../stores/categories.js';
  import { authStore } from '../stores/auth.js';
  import { openAuth } from '../stores/ui.js';

  export let isOpen = false;
  export let onClose = () => {};
  export let onSuccess = () => {};

  let auth;
  authStore.subscribe(val => auth = val);

  // Active tab: 'stream' (Link extractor) | 'file' (Local upload)
  let activeTab = 'stream';

  // Category management
  let selectedCategory = '';
  let newCategoryName = '';
  let availableCategories = [];
  categoryStore.subscribe(val => availableCategories = val);

  // Stream extraction state
  let streamUrl = '';
  let customTitle = '';
  let customArtist = '';
  let showCustomFields = false;

  // Extraction live status
  let isExtracting = false;
  let extractStatus = 'idle'; // 'idle' | 'processing' | 'completed' | 'failed'
  let extractProgress = 0;
  let extractStage = '';
  let extractError = '';
  let extractedTrack = null;
  let pollInterval = null;

  // Local file upload state
  let localTitle = '';
  let localArtist = '';
  let audioFile = null;
  let isLocalUploading = false;
  let localError = '';

  function handleFileChange(e) {
    if (e.target.files && e.target.files[0]) {
      audioFile = e.target.files[0];
    }
  }

  function handleAddCustomCategory() {
    if (!newCategoryName.trim()) return;
    const added = categoryStore.addCategory(newCategoryName.trim());
    if (added) {
      selectedCategory = added;
    }
    newCategoryName = '';
  }

  function clearPolling() {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
  }

  onDestroy(() => {
    clearPolling();
  });

  // Start stream extraction
  async function handleStartExtraction() {
    if (!streamUrl.trim()) return;
    extractError = '';
    isExtracting = true;
    extractStatus = 'processing';
    extractProgress = 6;
    extractStage = 'Connecting to stream source...';
    extractedTrack = null;

    try {
      const res = await tracksApi.extractStream({
        url: streamUrl.trim(),
        mood_tag: selectedCategory.trim() || 'General',
        custom_title: customTitle.trim(),
        custom_artist: customArtist.trim()
      });

      const jobId = res.job_id;
      if (!jobId) throw new Error('No job ID returned from server');

      // Poll progress every 500ms
      pollInterval = setInterval(async () => {
        try {
          const job = await tracksApi.getExtractJob(jobId);
          if (job) {
            extractProgress = Math.max(extractProgress, job.progress || 0);
            if (job.stage) extractStage = job.stage;

            if (job.status === 'completed') {
              clearPolling();
              extractStatus = 'completed';
              extractProgress = 100;
              extractedTrack = job.track;
              extractStage = 'Track successfully extracted!';
            } else if (job.status === 'failed') {
              clearPolling();
              extractStatus = 'failed';
              extractError = job.error || 'Failed to extract audio from stream';
            }
          }
        } catch (err) {
          // Soft poll error, continue retrying
          console.warn('Polling error:', err);
        }
      }, 500);

    } catch (err) {
      clearPolling();
      extractStatus = 'failed';
      extractError = err.message || 'Could not start stream extraction';
    }
  }

  function handleFinishAndClose() {
    clearPolling();
    onSuccess();
    handleCloseModal();
  }

  function handleResetExtraction() {
    clearPolling();
    isExtracting = false;
    extractStatus = 'idle';
    extractProgress = 0;
    extractStage = '';
    extractError = '';
    extractedTrack = null;
  }

  // Handle standard local file upload
  async function handleLocalUpload() {
    localError = '';
    isLocalUploading = true;

    try {
      const formData = new FormData();
      formData.append('title', localTitle.trim());
      formData.append('artist', localArtist.trim());
      formData.append('mood_tag', selectedCategory.trim() || 'General');

      if (audioFile) {
        formData.append('audio', audioFile);
      } else {
        throw new Error('Please select an audio file to upload');
      }

      await tracksApi.create(formData);
      onSuccess();
      handleCloseModal();
    } catch (err) {
      localError = err.message || 'Failed to upload audio file';
    } finally {
      isLocalUploading = false;
    }
  }

  function handleCloseModal() {
    clearPolling();
    isExtracting = false;
    extractStatus = 'idle';
    extractProgress = 0;
    streamUrl = '';
    customTitle = '';
    customArtist = '';
    localTitle = '';
    localArtist = '';
    audioFile = null;
    localError = '';
    extractError = '';
    extractedTrack = null;
    onClose();
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-calm-surface border border-calm-border rounded-lg shadow-xl max-w-md w-full p-6 sm:p-7 relative max-h-[90vh] overflow-y-auto text-calm-text">
      
      <!-- Close button -->
      <button 
        type="button"
        on:click={handleCloseModal}
        class="absolute top-5 right-5 text-neutral-400 hover:text-neutral-700 dark:hover:text-neutral-200 p-1 rounded-md transition-colors"
        aria-label="Close dialog"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>

      {#if !auth?.isAuthenticated}
        <div class="py-10 text-center">
          <div class="w-12 h-12 rounded-full bg-calm-ice text-calm-blue mx-auto flex items-center justify-center mb-4">
            <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
          </div>
          <h3 class="text-base sm:text-lg font-semibold text-calm-text mb-1">Login Diperlukan</h3>
          <p class="text-xs sm:text-sm text-calm-muted mb-6 max-w-sm mx-auto leading-relaxed">
            Anda harus masuk ke akun Anda terlebih dahulu untuk mengunggah atau mengekstrak lagu ke koleksi pribadi Anda.
          </p>
          <button
            type="button"
            on:click={() => { onClose(); openAuth(); }}
            class="px-5 py-2.5 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 transition-all shadow-sm"
          >
            Masuk / Buat Akun
          </button>
        </div>
      {:else}
        <!-- Modal Header -->
        <div class="mb-5">
          <h2 class="text-xl font-semibold tracking-tight text-calm-text">Add track</h2>
          <p class="text-xs text-calm-muted mt-0.5">Extract from a stream link or upload a local audio file.</p>
        </div>

      <!-- Mode Tab Bar (Explicit, intuitive, anti-slop) -->
      {#if !isExtracting}
        <div class="grid grid-cols-2 gap-1 p-1 bg-calm-mist/40 border border-calm-border rounded-md mb-5 text-xs">
          <button
            type="button"
            on:click={() => activeTab = 'stream'}
            class="py-2 px-3 rounded text-center transition-all font-medium flex items-center justify-center gap-1.5 {activeTab === 'stream' ? 'bg-calm-surface text-calm-text shadow-sm border border-calm-border/60' : 'text-calm-muted hover:text-calm-text'}"
          >
            <svg class="w-3.5 h-3.5 text-red-500" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19.615 3.184c-3.604-.246-11.631-.245-15.23 0-3.897.266-4.356 2.62-4.385 8.816.029 6.185.484 8.549 4.385 8.816 3.6.245 11.626.246 15.23 0 3.897-.266 4.356-2.62 4.385-8.816-.029-6.185-.484-8.549-4.385-8.816zm-10.615 12.816v-8l8 3.993-8 4.007z"/>
            </svg>
            <span>YouTube / Stream MP3</span>
          </button>
          <button
            type="button"
            on:click={() => activeTab = 'file'}
            class="py-2 px-3 rounded text-center transition-all font-medium flex items-center justify-center gap-1.5 {activeTab === 'file' ? 'bg-calm-surface text-calm-text shadow-sm border border-calm-border/60' : 'text-calm-muted hover:text-calm-text'}"
          >
            <svg class="w-3.5 h-3.5 text-calm-blue" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="17 8 12 3 7 8"></polyline>
              <line x1="12" y1="3" x2="12" y2="15"></line>
            </svg>
            <span>Upload File Lokal</span>
          </button>
        </div>
      {/if}

      <!-- CATEGORY SELECTION (Shared across modes when not extracting) -->
      {#if !isExtracting}
        <div class="mb-5">
          <div class="flex items-center justify-between mb-1.5">
            <label for="category-input-field" class="text-xs font-medium text-calm-text">Category</label>
            <span class="text-[11px] text-calm-muted">Optional</span>
          </div>

          <!-- Existing Category tags with remove button -->
          <div class="flex flex-wrap gap-1.5 mb-2.5">
            {#each availableCategories as cat}
              <div class="inline-flex items-center rounded border transition-colors text-xs {selectedCategory.toLowerCase() === cat.toLowerCase() ? 'bg-calm-text text-calm-bg border-calm-text font-medium' : 'bg-calm-surface text-calm-muted border-calm-border hover:border-calm-blue/60'}">
                <button 
                  type="button"
                  on:click={() => selectedCategory = cat}
                  class="px-2.5 py-1 text-xs"
                >
                  {cat}
                </button>
                <button 
                  type="button"
                  on:click|stopPropagation={() => {
                    categoryStore.removeCategory(cat);
                    if (selectedCategory.toLowerCase() === cat.toLowerCase()) selectedCategory = '';
                  }}
                  class="px-1.5 py-1 text-rose-500 hover:text-rose-700 hover:bg-rose-500/10 rounded-r transition-colors"
                  title="Remove category {cat}"
                  aria-label="Remove category {cat}"
                >
                  <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>
            {/each}
          </div>

          <!-- New category input -->
          <div class="flex gap-2">
            <input 
              id="category-input-field"
              type="text" 
              bind:value={newCategoryName} 
              placeholder="Create new category..."
              class="flex-1 px-3 py-1.5 rounded-md border border-calm-border bg-calm-mist/40 text-xs text-calm-text placeholder-calm-muted/60 focus:outline-none focus:ring-2 focus:ring-calm-blue focus:bg-calm-surface transition-colors"
            />
            <button 
              type="button"
              on:click={handleAddCustomCategory}
              disabled={!newCategoryName.trim()}
              class="px-3 py-1.5 rounded-md border border-calm-border bg-calm-surface text-calm-text hover:bg-calm-mist disabled:opacity-40 transition-colors text-xs font-medium"
            >
              Add
            </button>
          </div>
        </div>
      {/if}

      <!-- TAB 1: STREAM EXTRACTOR WITH REAL-TIME PROGRESS BAR -->
      {#if activeTab === 'stream'}
        {#if !isExtracting}
          <form on:submit|preventDefault={handleStartExtraction} class="space-y-4 text-xs">
            <!-- URL input -->
            <div>
              <label for="stream-url-input" class="block font-medium text-calm-text mb-1.5">Stream or video URL</label>
              <input 
                id="stream-url-input"
                type="url" 
                bind:value={streamUrl} 
                required 
                placeholder="e.g. https://www.youtube.com/watch?v=... or SoundCloud / Bandcamp"
                class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/40 text-xs text-calm-text placeholder-calm-muted/60 focus:outline-none focus:ring-2 focus:ring-calm-blue focus:bg-calm-surface transition-colors"
              />
              <p class="text-[11px] text-calm-muted mt-1">Audio will be extracted as 320k MP3 and saved directly to your private SeaweedFS storage.</p>
            </div>

            <!-- Optional customization accordion -->
            <div class="border border-calm-border/60 rounded-md p-3 bg-calm-mist/20">
              <button
                type="button"
                on:click={() => showCustomFields = !showCustomFields}
                class="flex items-center justify-between w-full text-left text-calm-muted hover:text-calm-text transition-colors"
              >
                <span class="font-medium text-[11px]">Customize title & artist (optional)</span>
                <svg class="w-3.5 h-3.5 transition-transform {showCustomFields ? 'rotate-180' : ''}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </button>

              {#if showCustomFields}
                <div class="mt-3 space-y-2.5 pt-2 border-t border-calm-border/40">
                  <div>
                    <label for="custom-title" class="block text-[11px] text-calm-muted mb-1">Track title</label>
                    <input 
                      id="custom-title"
                      type="text" 
                      bind:value={customTitle} 
                      placeholder="Leave empty to auto-detect"
                      class="w-full px-2.5 py-1.5 rounded border border-calm-border bg-calm-surface text-xs text-calm-text focus:outline-none focus:ring-1 focus:ring-calm-blue"
                    />
                  </div>
                  <div>
                    <label for="custom-artist" class="block text-[11px] text-calm-muted mb-1">Artist</label>
                    <input 
                      id="custom-artist"
                      type="text" 
                      bind:value={customArtist} 
                      placeholder="Leave empty to auto-detect"
                      class="w-full px-2.5 py-1.5 rounded border border-calm-border bg-calm-surface text-xs text-calm-text focus:outline-none focus:ring-1 focus:ring-calm-blue"
                    />
                  </div>
                </div>
              {/if}
            </div>

            <!-- Submit Button -->
            <div class="pt-2">
              <button 
                type="submit" 
                disabled={!streamUrl.trim()}
                class="w-full h-11 rounded-md bg-calm-text text-calm-bg text-sm font-medium hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-50 shadow-sm flex items-center justify-center gap-2"
              >
                <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                  <polyline points="7 10 12 15 17 10"></polyline>
                  <line x1="12" y1="15" x2="12" y2="3"></line>
                </svg>
                Extract and import track
              </button>
            </div>
          </form>

        {:else}
          <!-- ACTIVE EXTRACTION PROGRESS VIEW (As requested by user!) -->
          <div class="space-y-4 py-2 text-xs">
            
            <!-- Heading state -->
            <div class="flex items-center justify-between">
              <span class="font-medium text-calm-text text-sm">
                {#if extractStatus === 'processing'}
                  Extracting audio stream
                {:else if extractStatus === 'completed'}
                  Extraction complete
                {:else if extractStatus === 'failed'}
                  Extraction failed
                {/if}
              </span>
              <span class="font-medium text-xs {extractStatus === 'completed' ? 'text-emerald-500' : 'text-calm-muted'}">
                {extractProgress}%
              </span>
            </div>

            <!-- Animated Progress Bar -->
            <div class="w-full bg-calm-mist/60 rounded-full h-2.5 overflow-hidden border border-calm-border/40">
              <div 
                class="h-full rounded-full transition-all duration-300 ease-out {extractStatus === 'completed' ? 'bg-emerald-500' : extractStatus === 'failed' ? 'bg-rose-500' : 'bg-calm-text'}"
                style="width: {extractProgress}%;"
              ></div>
            </div>

            <!-- Current Stage Description -->
            <div class="flex items-center gap-2 text-xs text-calm-muted">
              {#if extractStatus === 'processing'}
                <svg class="animate-spin w-3.5 h-3.5 text-calm-blue flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <circle cx="12" cy="12" r="10" stroke-width="3" stroke-dasharray="32" stroke-linecap="round"></circle>
                </svg>
                <span class="truncate">{extractStage || 'Processing stream...'}</span>
              {:else if extractStatus === 'completed'}
                <svg class="w-4 h-4 text-emerald-500 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
                <span class="text-emerald-600 dark:text-emerald-400 font-medium">Lagu berhasil diekstrak dan disimpan ke storage!</span>
              {:else if extractStatus === 'failed'}
                <svg class="w-4 h-4 text-rose-500 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="15" y1="9" x2="9" y2="15"></line>
                  <line x1="9" y1="9" x2="15" y2="15"></line>
                </svg>
                <span class="text-rose-500">{extractError || 'Gagal mengekstrak stream.'}</span>
              {/if}
            </div>

            <!-- SUCCESS BANNER & PREVIEW CARD (User specifically requested: "dikasih tau baru masuk lagunya") -->
            {#if extractStatus === 'completed'}
              <div class="p-3.5 rounded-lg border border-emerald-500/20 bg-emerald-500/5 space-y-3 mt-3 animate-fade-in">
                <div class="flex items-center gap-2 text-xs text-emerald-600 dark:text-emerald-400 font-medium">
                  <span>Sukses! Lagu siap diputar di perpustakaan Anda.</span>
                </div>

                {#if extractedTrack}
                  <div class="flex items-center gap-3 p-2 rounded bg-calm-surface border border-calm-border">
                    <div class="w-11 h-11 rounded bg-calm-mist overflow-hidden flex-shrink-0 flex items-center justify-center">
                      <svg class="w-5 h-5 text-calm-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 18V5l12-2v13"></path>
                        <circle cx="6" cy="18" r="3"></circle>
                        <circle cx="18" cy="16" r="3"></circle>
                      </svg>
                    </div>
                    <div class="flex-1 min-w-0">
                      <div class="font-medium text-xs text-calm-text truncate">{extractedTrack.title}</div>
                      <div class="text-[11px] text-calm-muted truncate">{extractedTrack.artist}</div>
                    </div>
                    <span class="text-[10px] px-2 py-0.5 rounded border border-calm-border text-calm-muted">
                      {extractedTrack.mood_tag || 'Track'}
                    </span>
                  </div>
                {/if}

                <div class="pt-1">
                  <button
                    type="button"
                    on:click={handleFinishAndClose}
                    class="w-full h-10 rounded-md bg-calm-text text-calm-bg text-xs font-medium hover:opacity-90 active:scale-[0.99] transition-all shadow-sm"
                  >
                    Masuk ke perpustakaan (Selesai)
                  </button>
                </div>
              </div>
            {/if}

            <!-- ERROR STATE ACTIONS -->
            {#if extractStatus === 'failed'}
              <div class="pt-2 flex gap-2">
                <button
                  type="button"
                  on:click={handleResetExtraction}
                  class="flex-1 h-10 rounded-md border border-calm-border bg-calm-surface text-calm-text text-xs font-medium hover:bg-calm-mist transition-colors"
                >
                  Coba lagi
                </button>
                <button
                  type="button"
                  on:click={handleCloseModal}
                  class="px-4 h-10 rounded-md text-xs text-calm-muted hover:text-calm-text transition-colors"
                >
                  Tutup
                </button>
              </div>
            {/if}
          </div>
        {/if}

      <!-- TAB 2: LOCAL FILE UPLOAD -->
      {:else if activeTab === 'file'}
        <form on:submit|preventDefault={handleLocalUpload} class="space-y-4 text-xs">
          {#if localError}
            <div class="p-2.5 rounded-md bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs">
              {localError}
            </div>
          {/if}

          <!-- Title -->
          <div>
            <label for="local-track-title" class="block font-medium text-calm-text mb-1.5">Track title</label>
            <input 
              id="local-track-title"
              type="text" 
              bind:value={localTitle} 
              required 
              placeholder="e.g. Ref:rain"
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/40 text-xs text-calm-text placeholder-calm-muted/60 focus:outline-none focus:ring-2 focus:ring-calm-blue focus:bg-calm-surface transition-colors"
            />
          </div>

          <!-- Artist -->
          <div>
            <label for="local-track-artist" class="block font-medium text-calm-text mb-1.5">Artist</label>
            <input 
              id="local-track-artist"
              type="text" 
              bind:value={localArtist} 
              required 
              placeholder="e.g. Aimer"
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/40 text-xs text-calm-text placeholder-calm-muted/60 focus:outline-none focus:ring-2 focus:ring-calm-blue focus:bg-calm-surface transition-colors"
            />
          </div>

          <!-- Audio File Input -->
          <div>
            <label for="local-track-audio" class="block font-medium text-calm-text mb-1.5">Local audio file</label>
            <input 
              id="local-track-audio"
              type="file" 
              accept="audio/*"
              required
              on:change={handleFileChange}
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/40 text-xs text-calm-text file:mr-3 file:py-1 file:px-2.5 file:rounded file:border-0 file:text-xs file:font-medium file:bg-calm-text file:text-calm-bg hover:file:opacity-90 transition-all cursor-pointer"
            />
            <p class="text-[11px] text-calm-muted mt-1">MP3, WAV, FLAC, or AAC up to 32MB.</p>
          </div>

          <!-- Submit Button -->
          <div class="pt-2">
            <button 
              type="submit" 
              disabled={isLocalUploading || !audioFile}
              class="w-full h-11 rounded-md bg-calm-text text-calm-bg text-sm font-medium hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-50 shadow-sm flex items-center justify-center gap-2"
            >
              {#if isLocalUploading}
                <svg class="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <circle cx="12" cy="12" r="10" stroke-width="3" stroke-dasharray="32" stroke-linecap="round"></circle>
                </svg>
                Uploading to storage...
              {:else}
                Upload track
              {/if}
            </button>
          </div>
        </form>
      {/if}
    {/if}

    </div>
  </div>
{/if}
