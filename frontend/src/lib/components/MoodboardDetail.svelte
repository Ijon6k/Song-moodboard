<script>
  import { playerStore } from '../stores/player.js';
  import { authStore } from '../stores/auth.js';

  export let track;
  export let onBack = () => {};
  export let onOpenAddItem = () => {};

  let player;
  playerStore.subscribe(val => player = val);

  let auth;
  authStore.subscribe(val => auth = val);

  $: isCurrent = player?.currentTrack?.id === track.id;
  $: isPlaying = isCurrent && player?.isPlaying;

  function togglePlay() {
    playerStore.playTrack(track);
  }
</script>

<div class="moodboard-view animate-fade-in">
  <!-- Top Navigation & Track Meta -->
  <div class="moodboard-header">
    <button class="btn-back" on:click={onBack}>
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="19" y1="12" x2="5" y2="12"></line>
        <polyline points="12 19 5 12 12 5"></polyline>
      </svg>
      Kembali ke Katalog
    </button>

    <div class="header-actions">
      <button class="btn btn-secondary" on:click={onOpenAddItem}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        Tambah Serpihan Mood
      </button>
    </div>
  </div>

  <!-- Hero Track Banner -->
  <div class="track-banner shadow-calm">
    <div class="track-main-info">
      <div class="track-eyebrow">
        <span class="pulse-dot"></span>
        <span class="mood-tag">#{track.mood_tag || 'twilight'}</span>
      </div>
      <h2 class="track-heading">{track.title}</h2>
      <p class="track-sub">{track.artist}</p>
    </div>

    <div class="track-playback-ctrl">
      <button class="btn-play-hero {isPlaying ? 'playing' : ''}" on:click={togglePlay}>
        {#if isPlaying}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <rect x="6" y="4" width="4" height="16"></rect>
            <rect x="14" y="4" width="4" height="16"></rect>
          </svg>
          Jeda Nada
        {:else}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
          Putar Melodi
        {/if}
      </button>
    </div>
  </div>

  <!-- Moodboard Pinterest-style Masonry Grid -->
  <section class="grid-section">
    <div class="grid-title-row">
      <span class="section-label">Serpihan Visual & Renungan</span>
      <span class="section-count">{track.moodboard_items?.length || 0} fragmen</span>
    </div>

    {#if !track.moodboard_items || track.moodboard_items.length === 0}
      <div class="empty-state shadow-calm">
        <p class="empty-text">Belum ada serpihan visual atau kutipan di moodboard ini.</p>
        <button class="btn btn-secondary" on:click={onOpenAddItem}>+ Tambah Foto atau Kutipan</button>
      </div>
    {:else}
      <div class="masonry-grid">
        {#each track.moodboard_items as item}
          {#if item.item_type === 'image'}
            <div class="mood-card mood-card-image shadow-calm">
              <div class="img-frame">
                <img src={item.media_path} alt={item.content || "Mood photo"} loading="lazy" />
              </div>
              {#if item.content}
                <div class="card-caption">{item.content}</div>
              {/if}
            </div>
          {:else if item.item_type === 'quote'}
            <div class="mood-card mood-card-quote shadow-calm">
              <span class="quote-mark">“</span>
              <p class="quote-body">{item.content}</p>
              <div class="quote-footer">
                <span class="line"></span>
                <span class="quote-label">Catatan Suasana</span>
              </div>
            </div>
          {:else if item.item_type === 'palette'}
            <div class="mood-card mood-card-palette shadow-calm">
              <span class="palette-title">Palet Hening</span>
              <div class="swatch-row">
                {#each item.content.split(',') as color}
                  <div class="swatch-item">
                    <div class="swatch-pill" style="background-color: {color.trim()};"></div>
                    <span class="swatch-hex">{color.trim()}</span>
                  </div>
                {/each}
              </div>
            </div>
          {:else}
            <!-- Note / Memo -->
            <div class="mood-card mood-card-note shadow-calm">
              <div class="note-header">
                <span class="note-tag">Renungan</span>
              </div>
              <p class="note-body">{item.content}</p>
            </div>
          {/if}
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .moodboard-view {
    max-width: 1280px;
    margin: 0 auto;
    padding: 2rem 2rem 6rem;
  }

  .moodboard-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 2rem;
  }

  .btn-back {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: transparent;
    border: none;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--text-secondary);
    transition: color 0.15s ease;
  }

  .btn-back:hover {
    color: var(--text-primary);
  }

  .btn {
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    padding: 0.55rem 1.15rem;
    font-size: 0.875rem;
    font-weight: 500;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-secondary {
    background: var(--bg-ice);
    color: var(--accent-deep);
    border: 1px solid var(--border-blue);
  }

  .btn-secondary:hover {
    background: #deecf8;
    transform: translateY(-1px);
  }

  /* Track Banner */
  .track-banner {
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    padding: 2.2rem 2.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 3.5rem;
  }

  .track-eyebrow {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.6rem;
  }

  .pulse-dot {
    width: 7px;
    height: 7px;
    border-radius: 999px;
    background: var(--accent-blue);
    box-shadow: 0 0 8px var(--accent-blue);
  }

  .mood-tag {
    font-size: 0.8rem;
    text-transform: lowercase;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .track-heading {
    font-size: 1.85rem;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.25;
    letter-spacing: -0.02em;
    margin-bottom: 0.25rem;
  }

  .track-sub {
    font-size: 1rem;
    color: var(--text-secondary);
    font-weight: 400;
  }

  .btn-play-hero {
    display: inline-flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.8rem 1.5rem;
    border-radius: var(--radius-md);
    background: var(--text-primary);
    color: #ffffff;
    font-size: 0.95rem;
    font-weight: 500;
    border: none;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-play-hero:hover {
    background: #1e293b;
    transform: scale(1.02);
  }

  .btn-play-hero.playing {
    background: var(--accent-blue);
    color: #ffffff;
  }

  /* Grid section */
  .grid-section {
    margin-top: 2rem;
  }

  .grid-title-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.5rem;
    padding-bottom: 0.8rem;
    border-bottom: 1px solid var(--border-subtle);
  }

  .section-label {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.01em;
  }

  .section-count {
    font-size: 0.82rem;
    color: var(--text-muted);
  }

  /* Masonry Multi-Column CSS */
  .masonry-grid {
    column-count: 3;
    column-gap: 1.8rem;
  }

  @media (max-width: 1000px) {
    .masonry-grid {
      column-count: 2;
    }
  }

  @media (max-width: 650px) {
    .masonry-grid {
      column-count: 1;
    }
    .track-banner {
      flex-direction: column;
      align-items: flex-start;
      gap: 1.5rem;
      padding: 1.5rem;
    }
  }

  .mood-card {
    break-inside: avoid;
    margin-bottom: 1.8rem;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    overflow: hidden;
    transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;
  }

  .mood-card:hover {
    transform: translateY(-3px);
    border-color: var(--border-blue);
    box-shadow: var(--shadow-hover);
  }

  /* Image card */
  .mood-card-image .img-frame {
    width: 100%;
    overflow: hidden;
    background: var(--bg-mist);
  }

  .mood-card-image img {
    width: 100%;
    display: block;
    object-fit: cover;
    transition: transform 0.4s ease;
  }

  .mood-card-image:hover img {
    transform: scale(1.02);
  }

  .card-caption {
    padding: 0.9rem 1.1rem;
    font-size: 0.85rem;
    color: var(--text-secondary);
    line-height: 1.5;
  }

  /* Quote card */
  .mood-card-quote {
    background: var(--bg-ice);
    border-color: var(--border-blue);
    padding: 1.8rem;
    position: relative;
  }

  .quote-mark {
    font-size: 2.6rem;
    line-height: 1;
    color: var(--accent-blue);
    font-family: serif;
    display: block;
    margin-bottom: 0.3rem;
  }

  .quote-body {
    font-size: 0.98rem;
    font-style: italic;
    color: var(--text-primary);
    line-height: 1.6;
    margin-bottom: 1.2rem;
  }

  .quote-footer {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .quote-footer .line {
    width: 20px;
    height: 1px;
    background: var(--accent-blue);
  }

  .quote-label {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--accent-deep);
    font-weight: 500;
  }

  /* Palette card */
  .mood-card-palette {
    padding: 1.5rem;
  }

  .palette-title {
    display: block;
    font-size: 0.82rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }

  .swatch-row {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 0.5rem;
  }

  .swatch-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.35rem;
  }

  .swatch-pill {
    width: 100%;
    aspect-ratio: 1;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(0, 0, 0, 0.06);
  }

  .swatch-hex {
    font-size: 0.65rem;
    font-family: var(--font-mono);
    color: var(--text-muted);
  }

  /* Note card */
  .mood-card-note {
    padding: 1.5rem;
    background: #ffffff;
  }

  .note-header {
    margin-bottom: 0.8rem;
  }

  .note-tag {
    font-size: 0.72rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
    background: var(--bg-mist);
    padding: 0.2rem 0.5rem;
    border-radius: var(--radius-sm);
  }

  .note-body {
    font-size: 0.9rem;
    color: var(--text-secondary);
    line-height: 1.65;
  }

  .empty-state {
    padding: 3rem 2rem;
    text-align: center;
    background: var(--bg-card);
    border: 1px dashed var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .empty-text {
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }
</style>
