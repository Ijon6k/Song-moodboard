<script>
  import { openSearch } from "../stores/ui.js";

  export let onExplore = () => {};
  export let onUpload = () => {};

  let searchQuery = "";

  function handleSearchSubmit() {
    if (!searchQuery.trim()) return;
    openSearch(searchQuery.trim());
    searchQuery = "";
  }

  function handleSuggestionClick(tag) {
    openSearch(tag);
  }
</script>

<section
  class="relative w-full min-h-[100dvh] flex flex-col justify-between overflow-hidden bg-neutral-950"
>
  <!-- Full Viewport Background Photography (Original, no color filters) -->
  <div class="absolute inset-0 z-0">
    <img
      src="/images/hero.jpg"
      alt="Calm sky and atmosphere"
      class="w-full h-full object-cover object-center"
    />
  </div>

  <!-- Top spacer for fixed transparent navbar -->
  <div class="h-16 relative z-10"></div>

  <!-- Hero Content Center -->
  <div
    class="relative z-10 max-w-2xl w-full mx-auto px-6 py-12 text-center flex flex-col items-center justify-center my-auto"
  >
    <h1
      class="text-4xl sm:text-6xl font-semibold tracking-tight text-white leading-tight drop-shadow-[0_3px_12px_rgba(0,0,0,0.7)]"
    >
      Song Moodboard
    </h1>

    <p
      class="text-sm sm:text-lg text-white mt-3 max-w-md font-normal leading-relaxed drop-shadow-[0_1px_6px_rgba(0,0,0,0.7)]"
    >
      A quiet sanctuary for melodies, atmosphere, and visual memories.
    </p>

    <!-- Embedded Search Bar for Instant Music Exploration -->
    <form on:submit|preventDefault={handleSearchSubmit} class="w-full mt-8">
      <div
        class="relative flex items-center bg-black/45 backdrop-blur-xl border border-white/20 hover:border-white/35 focus-within:border-white/60 rounded-xl p-1.5 shadow-2xl transition-all"
      >
        <div class="pl-3.5 pr-2 text-neutral-400">
          <svg
            class="w-4 h-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </div>
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Search songs or artists across public music catalog..."
          class="flex-1 bg-transparent py-2 px-1 text-xs sm:text-sm text-white placeholder-neutral-300/70 focus:outline-none"
        />
        <button
          type="submit"
          class="px-4 py-2 rounded-lg bg-white text-neutral-950 text-xs sm:text-sm font-medium hover:bg-neutral-100 active:scale-[0.98] transition-all shadow-sm flex-shrink-0"
        >
          Search
        </button>
      </div>
    </form>

    <!-- Quick Suggestions -->
    <div
      class="flex flex-wrap items-center justify-center gap-1.5 mt-3 text-xs text-neutral-300"
    >
      <span class="text-neutral-400 text-[11px] mr-1">Suggestions:</span>
      {#each ["Aimer", "Joe Hisaishi", "Lofi Ambient", "Piano Relax", "Midnight Rain"] as tag}
        <button
          type="button"
          on:click={() => handleSuggestionClick(tag)}
          class="px-2.5 py-1 rounded-full bg-black/30 hover:bg-black/50 border border-white/15 text-neutral-200 hover:text-white text-[11px] transition-colors"
        >
          {tag}
        </button>
      {/each}
    </div>

    <!-- Secondary Actions Row -->
    <div
      class="flex flex-wrap items-center justify-center gap-3.5 mt-8 text-xs"
    >
      <button
        type="button"
        on:click={onUpload}
        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-white/15 hover:bg-white/25 border border-white/25 text-white font-medium backdrop-blur-sm transition-all shadow-sm active:scale-[0.98]"
      >
        <svg
          class="w-3.5 h-3.5 text-sky-300"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        <span>Upload local audio track</span>
      </button>

      <button
        type="button"
        on:click={onExplore}
        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg hover:bg-white/10 text-neutral-300 hover:text-white font-medium transition-colors"
      >
        <span>Browse catalog</span>
        <svg
          class="w-3.5 h-3.5"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </button>
    </div>
  </div>

  <!-- Bottom Bar with Scroll Indicator -->
  <div
    class="relative z-10 max-w-7xl w-full mx-auto px-6 pb-8 pt-4 flex items-center justify-center text-neutral-300"
  >
    <button
      type="button"
      on:click={onExplore}
      class="flex flex-col items-center gap-1.5 text-xs text-neutral-300 hover:text-white transition-colors cursor-pointer"
    >
      <span class="text-xs font-normal">Scroll to catalog</span>
      <div
        class="w-4 h-7 rounded-full border border-white/40 flex items-start justify-center p-1"
      >
        <div class="w-1 h-1.5 rounded-full bg-white animate-bounce"></div>
      </div>
    </button>
  </div>
</section>
