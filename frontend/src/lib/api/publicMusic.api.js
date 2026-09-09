/**
 * Public Music API Aggregator
 * Combines iTunes Search API (official previews, high-res artwork, mainstream releases)
 * and Audius Public API (full-length streaming tracks, ambient/lofi/indie releases).
 * Zero API keys required.
 */

export const publicMusicApi = {
  async search(query) {
    if (!query || !query.trim()) return [];
    const trimmed = query.trim();

    const [itunesRes, audiusRes] = await Promise.allSettled([
      fetch(`https://itunes.apple.com/search?term=${encodeURIComponent(trimmed)}&entity=song&limit=10`),
      fetch(`https://discoveryprovider.audius.co/v1/tracks/search?query=${encodeURIComponent(trimmed)}&app_name=Song Moodboard`)
    ]);

    const results = [];

    // Parse iTunes Results
    if (itunesRes.status === 'fulfilled' && itunesRes.value.ok) {
      try {
        const data = await itunesRes.value.json();
        if (data.results && Array.isArray(data.results)) {
          for (const item of data.results) {
            if (item.previewUrl) {
              const artworkHd = item.artworkUrl100 ? item.artworkUrl100.replace('100x100bb', '600x600bb') : null;
              results.push({
                id: `itunes-${item.trackId}`,
                rawId: item.trackId,
                title: item.trackName || 'Untitled',
                artist: item.artistName || 'Unknown Artist',
                audio_path: item.previewUrl,
                artwork: artworkHd || item.artworkUrl100,
                duration: item.trackTimeMillis ? Math.round(item.trackTimeMillis / 1000) : 30,
                mood_tag: item.primaryGenreName || 'General',
                source: 'iTunes Preview',
                isFull: false
              });
            }
          }
        }
      } catch (e) {
        console.warn('Failed parsing iTunes results:', e);
      }
    }

    // Parse Audius Results
    if (audiusRes.status === 'fulfilled' && audiusRes.value.ok) {
      try {
        const data = await audiusRes.value.json();
        if (data.data && Array.isArray(data.data)) {
          for (const item of data.data.slice(0, 8)) {
            const streamUrl = `https://api.audius.co/v1/tracks/${item.id}/stream?app_name=Song Moodboard`;
            const artwork = item.artwork?.['480x480'] || item.artwork?.['150x150'] || null;
            results.push({
              id: `audius-${item.id}`,
              rawId: item.id,
              title: item.title || 'Untitled',
              artist: item.user?.name || item.user?.handle || 'Unknown Artist',
              audio_path: streamUrl,
              artwork: artwork,
              duration: item.duration || 180,
              mood_tag: item.genre || item.mood || 'Ambient',
              source: 'Audius Full Track',
              isFull: true
            });
          }
        }
      } catch (e) {
        console.warn('Failed parsing Audius results:', e);
      }
    }

    return results;
  }
};
