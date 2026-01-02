import tailwindcss from "@tailwindcss/vite";
import { version } from "./package.json";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  ssr: false,
  css: ["~/assets/css/main.css"],
  vite: {
    plugins: [tailwindcss()],
  },
  modules: ["@nuxt/test-utils/module"],
  runtimeConfig: {
    public: {
      appVersion: version,
    },
  },
  nitro: {
    devProxy: {
      "/api": {
        target: "http://localhost:8080/api",
      },
    },
  },
});
