import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class",
  content: ["./src/**/*.tsx", "./src/**/*.ts", "./src/**/*.js"],

  theme: {
    extend: {
      fontFamily: {
        sans: ["Figtree", ...defaultTheme.fontFamily.sans],
      },
      colors: {
        shark: {
          50: "#F0F2F4",
          100: "#DEE2E8",
          200: "#BDC5D1",
          300: "#9CA8BA",
          400: "#7E8EA5",
          500: "#61728A",
          600: "#4A5769",
          700: "#323B48",
          800: "#1C2128",
          900: "#0F1115",
          950: "#060709",
        },
      },
      // Configuración de breakpoints personalizados
      screens: {
        xs: "475px", // Extra small devices (móviles pequeños)
        sm: "640px", // Small devices (móviles grandes)
        md: "768px", // Medium devices (tablets)
        lg: "1024px", // Large devices (laptops)
        xl: "1280px", // Extra large devices (desktops)
        "2xl": "1536px", // 2X large devices (large desktops)

        // Breakpoints específicos para mobile
        "mobile-s": "320px", // iPhone SE
        "mobile-m": "375px", // iPhone 12/13/14
        "mobile-l": "425px", // iPhone 14 Plus
        tablet: "768px", // iPad

        "max-sm": { max: "639px" }, // Hasta móviles
        "max-md": { max: "767px" }, // Hasta tablets pequeñas
        "max-lg": { max: "1023px" }, // Hasta tablets/laptops pequeños
      },
    },
  },
};
