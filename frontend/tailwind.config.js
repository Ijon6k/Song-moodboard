/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Poppins', '-apple-system', 'BlinkMacSystemFont', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
        serif: ['Lora', 'Georgia', 'serif'],
      },
      colors: {
        'calm-bg': 'var(--color-bg)',
        'calm-surface': 'var(--color-surface)',
        'calm-surface-elevated': 'var(--color-surface-elevated)',
        'calm-ice': 'var(--color-ice)',
        'calm-mist': 'var(--color-mist)',
        'calm-blue': {
          DEFAULT: '#8FB3D5',
          hover: '#799FC3',
          subtle: 'var(--color-blue-subtle)',
          deep: 'var(--color-blue-deep)',
        },
        'calm-text': 'var(--color-text)',
        'calm-muted': 'var(--color-muted)',
        'calm-border': 'var(--color-border)',
      },
      borderRadius: {
        DEFAULT: '8px',
        md: '8px',
        sm: '6px',
        pill: '9999px',
      },
      boxShadow: {
        calm: 'var(--shadow-calm)',
        hover: 'var(--shadow-hover)',
      },
    },
  },
  plugins: [],
};
