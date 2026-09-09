import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

const DEFAULT_CATEGORIES = [
  'Twilight',
  'Dawn',
  'Rain',
  'Midnight',
  'Acoustic',
  'Summer Sky'
];

function loadStoredCategories() {
  if (!browser) return DEFAULT_CATEGORIES;
  try {
    const raw = localStorage.getItem('songmoodboard_custom_categories');
    if (raw) {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed) && parsed.length > 0) {
        return parsed;
      }
    }
  } catch (e) {
    console.error('Failed to parse stored categories', e);
  }
  return DEFAULT_CATEGORIES;
}

function createCategoryStore() {
  const { subscribe, update, set } = writable(loadStoredCategories());

  return {
    subscribe,
    addCategory: (rawName) => {
      if (!rawName) return null;
      const trimmed = rawName.trim();
      if (!trimmed) return null;

      // Format neatly with capitalized first letter if lowercase
      const formatted = trimmed.charAt(0).toUpperCase() + trimmed.slice(1);

      update(list => {
        const exists = list.some(c => c.toLowerCase() === formatted.toLowerCase());
        if (exists) return list;
        const updated = [...list, formatted];
        if (browser) {
          try {
            localStorage.setItem('songmoodboard_custom_categories', JSON.stringify(updated));
          } catch (e) {}
        }
        return updated;
      });

      return formatted;
    },
    removeCategory: (name) => {
      update(list => {
        const updated = list.filter(c => c.toLowerCase() !== name.toLowerCase());
        if (browser) {
          try {
            localStorage.setItem('songmoodboard_custom_categories', JSON.stringify(updated));
          } catch (e) {}
        }
        return updated;
      });
    },
    reset: () => {
      if (browser) {
        localStorage.setItem('songmoodboard_custom_categories', JSON.stringify(DEFAULT_CATEGORIES));
      }
      set(DEFAULT_CATEGORIES);
    }
  };
}

export const categoryStore = createCategoryStore();
