import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createThemeStore() {
  const initialTheme = browser 
    ? (localStorage.getItem('songmoodboard_theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
    : 'light';

  const { subscribe, set, update } = writable(initialTheme);

  return {
    subscribe,
    toggleTheme: () => {
      update(current => {
        const next = current === 'dark' ? 'light' : 'dark';
        if (browser) {
          localStorage.setItem('songmoodboard_theme', next);
          if (next === 'dark') {
            document.documentElement.classList.add('dark');
          } else {
            document.documentElement.classList.remove('dark');
          }
        }
        return next;
      });
    },
    setTheme: (theme) => {
      if (browser) {
        localStorage.setItem('songmoodboard_theme', theme);
        if (theme === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
      }
      set(theme);
    },
    init: () => {
      if (browser) {
        const stored = localStorage.getItem('songmoodboard_theme');
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const active = stored || (prefersDark ? 'dark' : 'light');
        if (active === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
        set(active);
      }
    }
  };
}

export const themeStore = createThemeStore();
