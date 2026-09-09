<script>
  import { createQuery } from '@tanstack/svelte-query';
  import { tracksApi } from '$lib/api/tracks.api.js';
  import { queryKeys } from '$lib/queries/keys.js';
  import { categoryStore } from '$lib/stores/categories.js';
  import { authStore } from '$lib/stores/auth.js';
  import { openUpload, openAuth, openSearch } from '$lib/stores/ui.js';
  import Hero from '$lib/components/Hero.svelte';
  import TrackCard from '$lib/components/TrackCard.svelte';

  let auth;
  authStore.subscribe(val => auth = val);

  $: tracksQuery = createQuery({
    queryKey: queryKeys.tracks.all,
    queryFn: tracksApi.getAll,
    enabled: !!auth?.isAuthenticated,
  });

  let selectedCategory = 'all';
  let isAddCategoryOpen = false;
  let newCategoryInput = '';
  let categoryError = '';

  let customCategories = [];
  categoryStore.subscribe(val => customCategories = val);

  $: tracks = (auth?.isAuthenticated && $tracksQuery?.data) ? $tracksQuery.data : [];
  
  // Combine track categories with user custom categories
  $: allCategories = (() => {
    const fromTracks = tracks.map(t => t.mood_tag).filter(Boolean);
    const set = new Set(['all', ...customCategories, ...fromTracks]);
    return Array.from(set);
  })();

  $: filteredTracks = selectedCategory === 'all' 
    ? tracks 
    : tracks.filter(t => (t.mood_tag || 'General').toLowerCase() === selectedCategory.toLowerCase());

  function scrollToCatalog() {
    const el = document.getElementById('catalog');
    if (el) el.scrollIntoView({ behavior: 'smooth' });
  }

  function handleUploadClick() {
    if (!auth?.isAuthenticated) {
      openAuth();
    } else {
      openUpload();
    }
  }

  function handleUploadSuccess() {
    if ($tracksQuery?.refetch) {
      $tracksQuery.refetch();
    }
  }

  function handleCreateCategory() {
    categoryError = '';
    const trimmed = newCategoryInput.trim();
    if (!trimmed) {
      categoryError = 'Please enter a category name';
      return;
    }

    const created = categoryStore.addCategory(trimmed);
    if (created) {
      selectedCategory = created;
    }
    newCategoryInput = '';
    isAddCategoryOpen = false;
  }
</script>

<div>
  <!-- Full-viewport Hero with Calm Sky Photography & Serene Copywriting -->
  <Hero 
    onExplore={scrollToCatalog} 
    onUpload={handleUploadClick}
  />

  <!-- Main Soundscape Catalog Section -->
  <section id="catalog" class="max-w-7xl mx-auto px-6 py-16">
    <!-- Section Header & Filter Pills -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-5 mb-10 pb-5 border-b border-calm-border">
      <div>
        <div class="text-xs text-calm-muted font-medium mb-1">Catalog</div>
        <h2 class="text-xl sm:text-2xl font-semibold tracking-tight text-calm-text">Sound Collection</h2>
        <p class="text-xs sm:text-sm text-calm-muted mt-1">
          {#if auth?.isAuthenticated}
            Personal sanctuary for your soundscapes, moodboard fragments, and audio memories.
          {:else}
            Sign in to access your personal soundscapes, visual moodboards, and audio library.
          {/if}
        </p>
      </div>

      {#if auth?.isAuthenticated}
        <!-- Categories & Add Category Action -->
        <div class="flex items-center flex-wrap gap-2">
          <div class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 max-w-full">
            {#each allCategories as cat}
              <div class="inline-flex items-center rounded-sm border transition-all {selectedCategory.toLowerCase() === cat.toLowerCase() ? 'bg-calm-text text-calm-bg border-calm-text shadow-sm' : 'bg-calm-surface border-calm-border text-calm-muted hover:text-calm-text hover:border-calm-blue/50'}">
                <button 
                  type="button"
                  on:click={() => selectedCategory = cat}
                  class="px-3 py-1.5 text-xs font-medium whitespace-nowrap"
                >
                  {cat === 'all' ? 'All Atmospheres' : `#${cat}`}
                </button>
                {#if cat !== 'all' && customCategories.some(c => c.toLowerCase() === cat.toLowerCase())}
                  <button 
                    type="button"
                    on:click|stopPropagation={() => {
                      categoryStore.removeCategory(cat);
                      if (selectedCategory.toLowerCase() === cat.toLowerCase()) selectedCategory = 'all';
                    }}
                    class="pr-2 pl-0.5 py-1 text-rose-500 hover:text-rose-700 transition-colors"
                    title="Delete category {cat}"
                    aria-label="Delete #{cat}"
                  >
                    <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                {/if}
              </div>
            {/each}
          </div>

          <!-- Add Custom Category Trigger -->
          <button 
            type="button"
            on:click={() => isAddCategoryOpen = true}
            class="inline-flex items-center gap-1 px-3 py-1.5 rounded-sm text-xs font-medium border border-dashed border-calm-border bg-calm-surface text-calm-muted hover:text-calm-text hover:border-calm-blue transition-all"
            title="Create a new custom category"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
            <span>Category</span>
          </button>
        </div>
      {/if}
    </div>

    <!-- State rendering -->
    {#if !auth?.isAuthenticated}
      <!-- Unauthenticated Sanctuary Gate -->
      <div class="p-10 sm:p-14 text-center bg-calm-surface border border-calm-border rounded-xl shadow-sm backdrop-blur-sm animate-fade-in">
        <div class="w-12 h-12 rounded-full bg-calm-ice text-calm-blue mx-auto flex items-center justify-center mb-4">
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
        </div>
        <h3 class="text-lg sm:text-xl font-semibold text-calm-text mb-2">Personal Sound Sanctuary</h3>
        <p class="text-xs sm:text-sm text-calm-muted max-w-md mx-auto mb-6 leading-relaxed">
          Setiap soundscape dan moodboard bersifat privat dan unik bagi setiap akun. Silakan masuk atau daftar untuk melihat, mengunggah, dan mengelola koleksi Anda.
        </p>
        <div class="flex flex-wrap items-center justify-center gap-3">
          <button 
            type="button"
            on:click={openAuth}
            class="px-5 py-2.5 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 transition-all shadow-sm flex items-center gap-2"
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
              <polyline points="10 17 15 12 10 7"></polyline>
              <line x1="15" y1="12" x2="3" y2="12"></line>
            </svg>
            <span>Sign In / Register</span>
          </button>
          <button 
            type="button"
            on:click={() => openSearch()}
            class="px-4 py-2.5 rounded-md border border-calm-border bg-calm-surface text-calm-text text-xs sm:text-sm font-medium hover:bg-calm-mist transition-colors shadow-sm"
          >
            Search Public Music
          </button>
        </div>
      </div>
    {:else if $tracksQuery?.isLoading}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each Array(6) as _}
          <div class="h-48 bg-calm-surface border border-calm-border rounded-md p-6 animate-pulse flex flex-col justify-between">
            <div class="w-16 h-4 bg-calm-mist rounded-sm"></div>
            <div>
              <div class="w-3/4 h-5 bg-calm-mist rounded mb-2"></div>
              <div class="w-1/2 h-3 bg-calm-mist rounded"></div>
            </div>
            <div class="w-full h-3 bg-calm-mist rounded"></div>
          </div>
        {/each}
      </div>
    {:else if $tracksQuery?.isError}
      <div class="p-8 text-center bg-calm-surface border border-rose-500/30 rounded-md text-rose-600 dark:text-rose-400 text-xs sm:text-sm">
        Failed to load track catalog: {$tracksQuery.error?.message}
      </div>
    {:else if filteredTracks.length === 0}
      {#if tracks.length === 0}
        <!-- Logged in but empty library -->
        <div class="p-14 sm:p-16 text-center bg-calm-surface border border-dashed border-calm-border rounded-xl animate-fade-in">
          <div class="w-12 h-12 rounded-full bg-calm-ice text-calm-blue mx-auto flex items-center justify-center mb-3">
            <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <path d="M9 18V5l12-2v13"></path>
              <circle cx="6" cy="18" r="3"></circle>
              <circle cx="18" cy="16" r="3"></circle>
            </svg>
          </div>
          <h3 class="text-base sm:text-lg font-semibold text-calm-text mb-1">Koleksi lagu Anda masih kosong</h3>
          <p class="text-xs sm:text-sm text-calm-muted mb-6 max-w-md mx-auto leading-relaxed">
            Belum ada trek audio di akun Anda. Unggah file audio lokal atau ekstrak dari tautan streaming untuk mulai menyusun moodboard pribadi Anda.
          </p>
          <div class="flex flex-wrap items-center justify-center gap-3">
            <button 
              type="button"
              on:click={openUpload}
              class="inline-flex items-center gap-2 px-4 py-2.5 rounded-md bg-calm-text text-calm-bg text-xs sm:text-sm font-medium hover:opacity-90 transition-all shadow-sm"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
              <span>Upload Track</span>
            </button>
            <button 
              type="button"
              on:click={() => openSearch()}
              class="px-4 py-2.5 rounded-md border border-calm-border bg-calm-surface text-calm-text text-xs sm:text-sm font-medium hover:bg-calm-mist transition-colors shadow-sm"
            >
              Search Public Music
            </button>
          </div>
        </div>
      {:else}
        <!-- Filter resulted in 0 tracks -->
        <div class="p-16 text-center bg-calm-surface border border-dashed border-calm-border rounded-md">
          <div class="w-10 h-10 rounded-full bg-calm-ice text-calm-blue-deep mx-auto flex items-center justify-center mb-3">
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
          </div>
          <h3 class="text-sm font-semibold text-calm-text mb-1">No soundscapes found in #{selectedCategory}</h3>
          <p class="text-xs text-calm-muted mb-4 max-w-sm mx-auto">Upload an audio track under this category or choose another atmosphere.</p>
          <button 
            type="button"
            on:click={openUpload}
            class="px-4 py-2 rounded-md bg-calm-text text-calm-bg text-xs font-medium hover:opacity-90 transition-all shadow-sm"
          >
            Upload Track in this Category
          </button>
        </div>
      {/if}
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each filteredTracks as track (track.id)}
          <TrackCard {track} on:deleted={() => $tracksQuery.refetch()} />
        {/each}
      </div>
    {/if}
  </section>

  <!-- Add Category Modal -->
  {#if isAddCategoryOpen}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-sm animate-fade-in">
      <div class="bg-calm-surface rounded-md border border-calm-border shadow-hover max-w-sm w-full p-6 relative">
        <button 
          type="button"
          on:click={() => isAddCategoryOpen = false}
          class="absolute top-4 right-4 text-calm-muted hover:text-calm-text p-1 rounded-sm hover:bg-calm-mist transition-colors"
          aria-label="Close"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>

        <div class="mb-4">
          <h3 class="text-base font-semibold text-calm-text">Create Category</h3>
          <p class="text-xs text-calm-muted mt-0.5">Add a new atmosphere to organize your sounds and moodboards.</p>
        </div>

        {#if categoryError}
          <div class="mb-3 p-2 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-400 text-xs">
            {categoryError}
          </div>
        {/if}

        <form on:submit|preventDefault={handleCreateCategory} class="space-y-3">
          <div class="relative">
            <span class="absolute left-3 top-2.5 text-xs text-calm-muted font-mono">#</span>
            <input 
              type="text" 
              bind:value={newCategoryInput} 
              placeholder="e.g. Rainy Day, Stargazing, Piano"
              class="w-full pl-7 pr-3 py-2 rounded-md border border-calm-border bg-calm-mist/50 text-xs text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-2">
            <button 
              type="button"
              on:click={() => isAddCategoryOpen = false}
              class="px-3 py-1.5 rounded-md border border-calm-border text-calm-muted hover:text-calm-text text-xs font-medium transition-colors"
            >
              Cancel
            </button>
            <button 
              type="submit"
              class="px-4 py-1.5 rounded-md bg-calm-text text-calm-bg text-xs font-medium hover:opacity-90 transition-all shadow-sm"
            >
              Create Category
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</div>
