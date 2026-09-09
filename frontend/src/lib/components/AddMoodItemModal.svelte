<script>
  import { tracksApi } from '../api/tracks.api.js';

  export let isOpen = false;
  export let trackId = '';
  export let onClose = () => {};
  export let onSuccess = () => {};

  let itemType = 'image'; // 'image', 'quote', 'note', 'palette'
  let content = '';
  let imageFile = null;
  let isLoading = false;
  let errorMessage = '';

  function handleFileChange(e) {
    if (e.target.files && e.target.files[0]) {
      imageFile = e.target.files[0];
    }
  }

  async function handleSubmit() {
    if (!trackId) return;
    errorMessage = '';
    isLoading = true;

    try {
      const formData = new FormData();
      formData.append('item_type', itemType);
      formData.append('content', content);
      if (itemType === 'image' && imageFile) {
        formData.append('image', imageFile);
      }

      await tracksApi.addMoodItem(trackId, formData);
      onSuccess();
      onClose();
      // Reset
      content = '';
      imageFile = null;
    } catch (err) {
      errorMessage = err.message || 'Failed to pin mood item';
    } finally {
      isLoading = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-calm-surface rounded-md border border-calm-border shadow-hover max-w-lg w-full p-6 sm:p-7 relative max-h-[90vh] overflow-y-auto">
      <button 
        type="button"
        on:click={onClose}
        class="absolute top-4 right-4 text-calm-muted hover:text-calm-text p-1.5 rounded-sm hover:bg-calm-mist transition-colors"
        aria-label="Close modal"
      >
        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>

      <div class="mb-5">
        <h3 class="text-lg font-semibold text-calm-text">Add Mood Fragment</h3>
        <p class="text-xs text-calm-muted mt-0.5">Attach photographs, lyrics, palettes, or personal notes to this track.</p>
      </div>

      {#if errorMessage}
        <div class="mb-4 p-3 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-400 text-xs">
          {errorMessage}
        </div>
      {/if}

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div>
          <label for="item-type-select" class="block text-xs font-medium text-calm-text mb-1.5">Fragment Type</label>
          <div id="item-type-select" class="grid grid-cols-4 gap-2">
            {#each [
              { id: 'image', label: 'Photo' },
              { id: 'quote', label: 'Quote' },
              { id: 'palette', label: 'Palette' },
              { id: 'note', label: 'Note' }
            ] as t}
              <button 
                type="button" 
                on:click={() => itemType = t.id}
                class="py-2 text-xs font-medium rounded-md border transition-all {itemType === t.id ? 'bg-calm-text text-calm-bg border-calm-text shadow-sm' : 'bg-calm-surface text-calm-muted border-calm-border hover:bg-calm-mist hover:text-calm-text'}"
              >
                {t.label}
              </button>
            {/each}
          </div>
        </div>

        {#if itemType === 'image'}
          <div>
            <label for="frag-image" class="block text-xs font-medium text-calm-text mb-1">Select Photograph (.jpg, .png, .webp)</label>
            <input 
              id="frag-image"
              type="file" 
              accept="image/*"
              required={!imageFile}
              on:change={handleFileChange}
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-xs text-calm-text file:mr-3 file:py-1.5 file:px-3 file:rounded-sm file:border-0 file:text-xs file:font-medium file:bg-calm-ice file:text-calm-blue-deep hover:file:bg-calm-blue hover:file:text-white transition-all cursor-pointer"
            />
          </div>
          <div>
            <label for="frag-img-caption" class="block text-xs font-medium text-calm-text mb-1">Photo Caption (Optional)</label>
            <input 
              id="frag-img-caption"
              type="text" 
              bind:value={content}
              placeholder="e.g. Misty morning sky over calm waters..."
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
            />
          </div>
        {:else if itemType === 'quote'}
          <div>
            <label for="frag-quote" class="block text-xs font-medium text-calm-text mb-1">Quote or Lyrics Fragment</label>
            <textarea 
              id="frag-quote"
              bind:value={content}
              required
              rows="3"
              placeholder="e.g. Even in the deepest twilight, faint stars persist with quiet grace..."
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors font-serif italic"
            ></textarea>
          </div>
        {:else if itemType === 'palette'}
          <div>
            <label for="frag-palette" class="block text-xs font-medium text-calm-text mb-1">Hex Color Palette (comma separated)</label>
            <input 
              id="frag-palette"
              type="text" 
              bind:value={content}
              required
              placeholder="#EBF3FA, #8FB3D5, #2E4B66, #F8FAFC"
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors font-mono"
            />
          </div>
        {:else}
          <div>
            <label for="frag-memo" class="block text-xs font-medium text-calm-text mb-1">Personal Note</label>
            <textarea 
              id="frag-memo"
              bind:value={content}
              required
              rows="3"
              placeholder="e.g. Best experienced at 5:00 AM while drinking warm tea by the open window..."
              class="w-full px-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
            ></textarea>
          </div>
        {/if}

        <div class="pt-2">
          <button 
            type="submit" 
            disabled={isLoading}
            class="w-full py-2.5 px-4 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-50 shadow-sm"
          >
            {#if isLoading}
              Pinning to Moodboard...
            {:else}
              Pin Fragment
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
