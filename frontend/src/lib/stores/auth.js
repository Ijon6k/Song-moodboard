import { writable } from 'svelte/store';

function createAuthStore() {
  const initialToken = typeof localStorage !== 'undefined' ? localStorage.getItem('sm_token') : null;
  const initialUser = typeof localStorage !== 'undefined' && localStorage.getItem('sm_user') 
    ? JSON.parse(localStorage.getItem('sm_user')) 
    : null;

  const { subscribe, set, update } = writable({
    token: initialToken,
    user: initialUser,
    isAuthenticated: !!initialToken,
    isLoading: false,
    error: null
  });

  return {
    subscribe,
    login: async (email, password) => {
      update(s => ({ ...s, isLoading: true, error: null }));
      try {
        const res = await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password })
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Login gagal');

        localStorage.setItem('sm_token', data.token);
        localStorage.setItem('sm_user', JSON.stringify(data.user));

        set({
          token: data.token,
          user: data.user,
          isAuthenticated: true,
          isLoading: false,
          error: null
        });
        return { success: true };
      } catch (err) {
        update(s => ({ ...s, isLoading: false, error: err.message }));
        return { success: false, error: err.message };
      }
    },

    register: async (username, email, password) => {
      update(s => ({ ...s, isLoading: true, error: null }));
      try {
        const res = await fetch('/api/auth/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, email, password })
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Pendaftaran gagal');

        localStorage.setItem('sm_token', data.token);
        localStorage.setItem('sm_user', JSON.stringify(data.user));

        set({
          token: data.token,
          user: data.user,
          isAuthenticated: true,
          isLoading: false,
          error: null
        });
        return { success: true };
      } catch (err) {
        update(s => ({ ...s, isLoading: false, error: err.message }));
        return { success: false, error: err.message };
      }
    },

    logout: async () => {
      const state = typeof localStorage !== 'undefined' ? localStorage.getItem('sm_token') : null;
      if (state) {
        try {
          await fetch('/api/auth/logout', {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${state}` }
          });
        } catch (e) {
          // ignore error on logout
        }
      }
      localStorage.removeItem('sm_token');
      localStorage.removeItem('sm_user');
      set({
        token: null,
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null
      });
    }
  };
}

export const authStore = createAuthStore();
