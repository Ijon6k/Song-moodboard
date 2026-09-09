<script>
  import '../app.css';
  import { browser } from '$app/environment';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
  import { themeStore } from '$lib/stores/theme.js';
  import Navbar from '$lib/components/Navbar.svelte';
  import AudioPlayer from '$lib/components/AudioPlayer.svelte';
  import AuthModal from '$lib/components/AuthModal.svelte';
  import UploadTrackModal from '$lib/components/UploadTrackModal.svelte';
  import DiscoverMusicModal from '$lib/components/DiscoverMusicModal.svelte';

  import { 
    isSearchOpen, 
    searchInitialQuery, 
    isUploadOpen, 
    isAuthOpen, 
    openSearch, 
    closeSearch, 
    openUpload, 
    closeUpload, 
    openAuth, 
    closeAuth 
  } from '$lib/stores/ui.js';

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        enabled: browser,
        staleTime: 1000 * 30,
        refetchOnWindowFocus: false,
      },
    },
  });

  onMount(() => {
    themeStore.init();
  });

  function handleUploadSuccess() {
    queryClient.invalidateQueries({ queryKey: ['tracks'] });
  }
</script>

<QueryClientProvider client={queryClient}>
  <div class="min-h-screen flex flex-col bg-calm-bg text-calm-text font-sans relative">
    <Navbar 
      onOpenAuth={() => openAuth()}
      onOpenUpload={() => openUpload()}
      onOpenSearch={() => openSearch()}
    />

    <main class="flex-1 pb-28 {$page.url.pathname === '/' ? '' : 'pt-16'}">
      <slot />
    </main>

    <!-- Global Floating Audio Player -->
    <AudioPlayer />

    <!-- Modals -->
    <AuthModal 
      isOpen={$isAuthOpen} 
      onClose={closeAuth} 
    />

    <UploadTrackModal 
      isOpen={$isUploadOpen} 
      onClose={closeUpload} 
      onSuccess={handleUploadSuccess}
    />

    <DiscoverMusicModal 
      isOpen={$isSearchOpen}
      initialQuery={$searchInitialQuery}
      onClose={closeSearch}
      onImportSuccess={handleUploadSuccess}
    />
  </div>
</QueryClientProvider>
