import { writable } from 'svelte/store';

function createPlayerStore() {
  const { subscribe, set, update } = writable({
    currentTrack: null,
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    volume: 0.8,
  });

  let audio = null;

  if (typeof window !== 'undefined') {
    audio = new Audio();
    audio.volume = 0.8;

    audio.addEventListener('timeupdate', () => {
      update(s => ({ ...s, currentTime: audio.currentTime }));
    });

    audio.addEventListener('loadedmetadata', () => {
      update(s => ({ ...s, duration: audio.duration }));
    });

    audio.addEventListener('ended', () => {
      update(s => ({ ...s, isPlaying: false, currentTime: 0 }));
    });

    audio.addEventListener('play', () => {
      update(s => ({ ...s, isPlaying: true }));
    });

    audio.addEventListener('pause', () => {
      update(s => ({ ...s, isPlaying: false }));
    });
  }

  return {
    subscribe,
    playTrack: (track) => {
      if (!audio) return;
      update(s => {
        if (s.currentTrack?.id === track.id) {
          if (s.isPlaying) {
            audio.pause();
          } else {
            audio.play().catch(console.error);
          }
          return s;
        }

        // New track
        const path = track.audio_path || '';
        audio.src = path.startsWith('/api') || path.startsWith('/demo-audio') || path.startsWith('http://') || path.startsWith('https://')
          ? path 
          : `/api/storage/${path}`;
        audio.play().catch(console.error);

        return {
          ...s,
          currentTrack: track,
          isPlaying: true,
          currentTime: 0
        };
      });
    },

    seekRelative: (delta) => {
      if (!audio) return;
      const target = Math.max(0, Math.min(audio.duration || 100, audio.currentTime + delta));
      audio.currentTime = target;
    },

    togglePlay: () => {
      if (!audio) return;
      if (audio.paused) {
        audio.play().catch(console.error);
      } else {
        audio.pause();
      }
    },

    seek: (timeInSeconds) => {
      if (!audio) return;
      audio.currentTime = timeInSeconds;
    },

    setVolume: (val) => {
      if (!audio) return;
      audio.volume = Math.max(0, Math.min(1, val));
      update(s => ({ ...s, volume: audio.volume }));
    }
  };
}

export const playerStore = createPlayerStore();
