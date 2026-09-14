/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        bg: {
          primary: '#0F1115',
          secondary: '#171A21',
          surface: '#1E222B',
          card: 'rgba(35, 40, 52, 0.75)',
        },
        glass: {
          DEFAULT: 'rgba(255, 255, 255, 0.08)',
          hover: 'rgba(255, 255, 255, 0.12)',
          border: 'rgba(255, 255, 255, 0.18)',
          separator: 'rgba(255, 255, 255, 0.08)',
        },
        text: {
          primary: '#FFFFFF',
          secondary: '#B8C0CC',
          muted: '#7D8592',
        },
        brand: {
          primary: '#0A84FF',
          hover: '#339CFF',
        },
        status: {
          success: '#30D158',
          warning: '#FFD60A',
          error: '#FF453A',
          info: '#64D2FF',
        },
      },
      fontFamily: {
        sans: [
          '"SF Pro Display"',
          '"SF Pro Text"',
          '-apple-system',
          'BlinkMacSystemFont',
          '"Segoe UI"',
          '"Helvetica Neue"',
          'Helvetica',
          'Arial',
          'sans-serif',
        ],
        mono: [
          '"SF Mono"',
          'ui-monospace',
          'SFMono-Regular',
          'Menlo',
          'Monaco',
          'Consolas',
          '"Liberation Mono"',
          '"Courier New"',
          'monospace',
        ],
      },
      borderRadius: {
        'glass': '18px',
        'glass-sm': '12px',
        'glass-lg': '22px',
        'liquid': '28px',
      },
      boxShadow: {
        'glass': '0 12px 40px rgba(0, 0, 0, 0.25)',
        'glass-hover': '0 16px 48px rgba(0, 0, 0, 0.35)',
        'glow-primary': '0 0 24px rgba(10, 132, 255, 0.35)',
        'specular': 'inset 0 1px 0 0 rgba(255, 255, 255, 0.5), inset 0 0 0 1px rgba(255, 255, 255, 0.04), 0 12px 40px rgba(0, 0, 0, 0.25)',
        'specular-soft': 'inset 0 1px 0 0 rgba(255, 255, 255, 0.22), inset 0 0 0 1px rgba(255, 255, 255, 0.03)',
      },
      transitionTimingFunction: {
        'liquid': 'cubic-bezier(0.22, 1, 0.36, 1)',
        'liquid-fast': 'cubic-bezier(0.3, 1.2, 0.4, 1)',
      },
      backdropBlur: {
        'glass': '36px',
        'glass-strong': '48px',
      },
    },
  },
  plugins: [],
}
