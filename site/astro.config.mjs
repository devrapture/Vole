import { defineConfig } from "astro/config";
import mdx from "@astrojs/mdx";
import tailwind from "@astrojs/tailwind";
import sitemap from "@astrojs/sitemap";
import expressiveCode from "astro-expressive-code";

export default defineConfig({
  site: "https://vole-clean.dev",
  integrations: [
    expressiveCode({
      themes: ["vitesse-dark"],
      styleOverrides: {
        borderRadius: "0.5rem",
        borderWidth: "1px",
      },
    }),
    mdx(),
    tailwind(),
    sitemap(),
  ],
});
