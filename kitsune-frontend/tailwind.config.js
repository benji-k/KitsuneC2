/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors:{
        'kc2-wine-purple': '#252340',
        'kc2-light-gray' : '#4F4D67',
        'kc2-dark-gray' : '#222132',
        'kc2-soap-pink' : '#E2AFEA',
        'kc2-dashboard-bg' : '#353446'
      },
    },
  },
  plugins: [],
}
