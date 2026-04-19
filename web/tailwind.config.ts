/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:55
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:55
 * @FilePath: /web/tailwind.config.ts
 * @Description: 
 * 
 */
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
