<script>
  import { authStore } from '../stores/auth.js';

  export let isOpen = false;
  export let onClose = () => {};

  let isRegister = false;
  let username = '';
  let email = '';
  let password = '';
  let errorMessage = '';
  let isLoading = false;

  async function handleSubmit() {
    errorMessage = '';
    isLoading = true;

    if (isRegister) {
      const res = await authStore.register(username, email, password);
      if (res.success) {
        onClose();
      } else {
        errorMessage = res.error;
      }
    } else {
      const res = await authStore.login(email, password);
      if (res.success) {
        onClose();
      } else {
        errorMessage = res.error;
      }
    }
    isLoading = false;
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-calm-surface rounded-md border border-calm-border shadow-hover max-w-md w-full p-6 sm:p-7 relative">
      <!-- Close button -->
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

      <!-- Tabs -->
      <div class="flex items-center gap-4 mb-6 border-b border-calm-border pb-3">
        <button 
          type="button"
          on:click={() => { isRegister = false; errorMessage = ''; }}
          class="text-sm font-semibold pb-1 transition-colors {!isRegister ? 'text-calm-text border-b-2 border-calm-blue' : 'text-calm-muted hover:text-calm-text'}"
        >
          Sign In
        </button>
        <button 
          type="button"
          on:click={() => { isRegister = true; errorMessage = ''; }}
          class="text-sm font-semibold pb-1 transition-colors {isRegister ? 'text-calm-text border-b-2 border-calm-blue' : 'text-calm-muted hover:text-calm-text'}"
        >
          Create Account
        </button>
      </div>

      {#if errorMessage}
        <div class="mb-4 p-3 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-400 text-xs">
          {errorMessage}
        </div>
      {/if}

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#if isRegister}
          <div>
            <label for="username" class="block text-xs font-medium text-calm-text mb-1">Username</label>
            <input 
              id="username"
              type="text" 
              bind:value={username} 
              required 
              placeholder="e.g. tranquil_listener"
              class="w-full px-3.5 py-2.5 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
            />
          </div>
        {/if}

        <div>
          <label for="email" class="block text-xs font-medium text-calm-text mb-1">Email Address</label>
          <input 
            id="email"
            type="email" 
            bind:value={email} 
            required 
            placeholder="calm@song-moodboard.io"
            class="w-full px-3.5 py-2.5 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
          />
        </div>

        <div>
          <label for="password" class="block text-xs font-medium text-calm-text mb-1">Password</label>
          <input 
            id="password"
            type="password" 
            bind:value={password} 
            required 
            minlength="6"
            placeholder="••••••••"
            class="w-full px-3.5 py-2.5 rounded-md border border-calm-border bg-calm-mist/50 text-sm text-calm-text placeholder-calm-muted/60 focus:outline-none focus:border-calm-blue focus:bg-calm-surface transition-colors"
          />
        </div>

        <button 
          type="submit" 
          disabled={isLoading}
          class="w-full py-2.5 px-4 rounded-md bg-calm-text text-calm-bg text-sm font-medium hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-50 mt-2 shadow-sm"
        >
          {#if isLoading}
            Authenticating...
          {:else}
            {isRegister ? 'Create Song Moodboard Account' : 'Sign In to Song Moodboard'}
          {/if}
        </button>
      </form>

      <div class="mt-5 text-center text-xs text-calm-muted">
        Demo Login: <span class="font-mono text-calm-text">calm@song-moodboard.io</span> / <span class="font-mono text-calm-text">silentpassword123</span>
      </div>
    </div>
  </div>
{/if}
