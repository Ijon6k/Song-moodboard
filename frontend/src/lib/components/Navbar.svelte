<script>
  import { page } from '$app/stores';
  import { authStore } from '../stores/auth.js';
  import { themeStore } from '../stores/theme.js';

  export let onOpenAuth = () => {};
  export let onOpenUpload = () => {};
  export let onOpenSearch = () => {};

  let auth;
  authStore.subscribe(val => auth = val);

  let theme;
  themeStore.subscribe(val => theme = val);

  let scrollY = 0;
  $: isHome = $page.url.pathname === '/';
  $: isScrolled = scrollY > 50 || !isHome;
</script>

<svelte:window bind:scrollY={scrollY} />

<header 
  class="fixed top-0 inset-x-0 z-40 h-16 transition-all duration-300 {
    isScrolled 
      ? 'bg-white/85 dark:bg-[#090D16]/90 backdrop-blur-md border-b border-calm-border shadow-sm' 
      : 'bg-transparent border-b border-transparent'
  }"
>
  <div class="max-w-7xl h-full mx-auto px-6 flex items-center justify-between">
    <!-- Brand -->
    <a href="/" class="flex items-center gap-2.5 group cursor-pointer">
      <span class="w-2 h-2 rounded-full transition-transform group-hover:scale-125 {isScrolled ? 'bg-calm-blue' : 'bg-white'}"></span>
      <span class="text-base font-semibold tracking-tight transition-colors {isScrolled ? 'text-calm-text' : 'text-white drop-shadow-sm'}">
        Song Moodboard
      </span>
    </a>

    <!-- Center Navigation (English) -->
    <nav class="hidden md:flex items-center gap-7 text-xs font-medium transition-colors {isScrolled ? 'text-calm-muted' : 'text-white/80'}">
      <a href="/" class="hover:{isScrolled ? 'text-calm-text' : 'text-white'} transition-colors">Collection</a>
      <a href="/about" class="hover:{isScrolled ? 'text-calm-text' : 'text-white'} transition-colors">Architecture & Specs</a>
    </nav>

    <!-- Right Actions -->
    <div class="flex items-center gap-2.5">
      <!-- Search Public Music Button (Only shows when scrolled, as hero already has search) -->
      {#if isScrolled}
        <button 
          type="button"
          on:click={onOpenSearch}
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md border border-calm-border bg-calm-surface text-calm-muted hover:text-calm-text hover:border-calm-blue transition-all shadow-sm animate-fade-in"
          title="Search public music API (iTunes & Audius)"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <span class="hidden sm:inline">Search Music</span>
        </button>
      {/if}

      <!-- Dark / Light Theme Toggle -->
      <button 
        type="button"
        on:click={() => themeStore.toggleTheme()}
        class="p-2 rounded-md transition-all shadow-sm {
          isScrolled 
            ? 'border border-calm-border bg-calm-surface text-calm-muted hover:text-calm-text hover:bg-calm-mist' 
            : 'border border-white/20 bg-black/30 backdrop-blur-sm text-white hover:bg-black/50'
        }"
        title={theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
        aria-label="Toggle theme mode"
      >
        {#if theme === 'dark'}
          <!-- Sun icon -->
          <svg class="w-4 h-4 text-amber-300 transition-transform rotate-0 hover:rotate-45" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="5"></circle>
            <line x1="12" y1="1" x2="12" y2="3"></line>
            <line x1="12" y1="21" x2="12" y2="23"></line>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
            <line x1="1" y1="12" x2="3" y2="12"></line>
            <line x1="21" y1="12" x2="23" y2="12"></line>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
          </svg>
        {:else}
          <!-- Moon icon -->
          <svg class="w-4 h-4 {isScrolled ? 'text-slate-600' : 'text-slate-200'} transition-transform hover:-rotate-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
          </svg>
        {/if}
      </button>

      {#if auth?.isAuthenticated}
        <button 
          on:click={onOpenUpload}
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-medium rounded-md transition-all shadow-sm {
            isScrolled 
              ? 'bg-calm-text text-calm-bg hover:opacity-90 hover:-translate-y-0.5' 
              : 'bg-white text-neutral-950 hover:bg-neutral-100 hover:-translate-y-0.5'
          }"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          Upload Track
        </button>

        <div class="flex items-center gap-2 py-1 px-2.5 rounded-md text-xs shadow-sm {
          isScrolled 
            ? 'bg-calm-surface border border-calm-border' 
            : 'bg-black/30 backdrop-blur-sm border border-white/20 text-white'
        }">
          <span class="w-5 h-5 rounded-sm bg-calm-blue text-white flex items-center justify-center font-semibold text-[10px]">
            {auth.user?.username ? auth.user.username[0].toUpperCase() : 'U'}
          </span>
          <span class="font-medium text-xs {isScrolled ? 'text-calm-text' : 'text-white'}">{auth.user?.username}</span>
          <button 
            on:click={() => authStore.logout()} 
            class="transition-colors p-0.5 ml-1 {isScrolled ? 'text-calm-muted hover:text-rose-500' : 'text-white/70 hover:text-rose-400'}" 
            title="Sign Out"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
          </button>
        </div>
      {:else}
        <a href="/about" class="md:hidden text-xs {isScrolled ? 'text-calm-muted hover:text-calm-text' : 'text-white/80 hover:text-white'}">Tech</a>
        <button 
          on:click={onOpenAuth}
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-medium rounded-md transition-all shadow-sm {
            isScrolled 
              ? 'bg-calm-text text-calm-bg hover:opacity-90 hover:-translate-y-0.5' 
              : 'bg-white text-neutral-950 hover:bg-neutral-100 hover:-translate-y-0.5'
          }"
        >
          Sign In
        </button>
      {/if}
    </div>
  </div>
</header>
