import { writable } from 'svelte/store';

export const isSearchOpen = writable(false);
export const searchInitialQuery = writable('');
export const isUploadOpen = writable(false);
export const isAuthOpen = writable(false);

export function openSearch(query = '') {
  searchInitialQuery.set(query);
  isSearchOpen.set(true);
}

export function closeSearch() {
  isSearchOpen.set(false);
}

export function openUpload() {
  isUploadOpen.set(true);
}

export function closeUpload() {
  isUploadOpen.set(false);
}

export function openAuth() {
  isAuthOpen.set(true);
}

export function closeAuth() {
  isAuthOpen.set(false);
}
