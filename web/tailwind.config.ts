import type { Config } from 'tailwindcss'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      spacing: {
        '3': '0.75rem',
      },
      gap: {
        '3': '0.75rem',
      },
      padding: {
        '3': '0.75rem',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
} satisfies Config
