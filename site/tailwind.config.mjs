import typography from "@tailwindcss/typography";

/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      colors: {
        bg: "#050505",
        surface: "#0d0d0f",
        border: "#262626",
        text: "#ffffff",
        muted: "#a3a3a3",
        accent: "#22c55e",
        warning: "#f59e0b",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      backgroundImage: {
        grid: `radial-gradient(circle, #262626 1px, transparent 1px)`,
      },
      backgroundSize: {
        grid: "24px 24px",
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: "none",
            fontSize: "1rem",
            lineHeight: "1.8",
            "--tw-prose-body": "#a3a3a3",
            "--tw-prose-headings": "#ffffff",
            "--tw-prose-links": "#22c55e",
            "--tw-prose-bold": "#ffffff",
            "--tw-prose-counters": "#a3a3a3",
            "--tw-prose-bullets": "#a3a3a3",
            "--tw-prose-hr": "#262626",
            "--tw-prose-quotes": "#a3a3a3",
            "--tw-prose-quote-borders": "#262626",
            "--tw-prose-code": "#ffffff",
            "--tw-prose-pre-code": "#ffffff",
            "--tw-prose-pre-bg": "#0d0d0f",
            "--tw-prose-th-borders": "#262626",
            "--tw-prose-td-borders": "#262626",
            h1: {
              fontSize: "2rem",
              fontWeight: 700,
              lineHeight: 1.25,
              marginTop: "1rem",
              marginBottom: "1rem",
            },
            h2: {
              fontSize: "1.5rem",
              fontWeight: 600,
              lineHeight: 1.3,
              marginTop: "4.5rem",
              marginBottom: "1.5rem",
              letterSpacing: "-0.01em",
            },
            h3: {
              fontSize: "1.25rem",
              fontWeight: 600,
              lineHeight: 1.4,
              marginTop: "3rem",
              marginBottom: "1.25rem",
              letterSpacing: "-0.01em",
            },
            p: {
              marginBottom: "1.5rem",
              lineHeight: "1.8",
            },
            "ul, ol": {
              marginBottom: "1.5rem",
              paddingLeft: "1.5rem",
            },
            li: {
              marginTop: "0.375rem",
              marginBottom: "0.375rem",
            },
            "li > ul, li > ol": {
              marginTop: "0.375rem",
              marginBottom: "0",
            },
            pre: {
              marginTop: "3.5rem",
              marginBottom: "3.5rem",
              borderRadius: "0.5rem",
              fontSize: "0.875rem",
              lineHeight: "1.6",
            },
            code: {
              borderRadius: "0.375rem",
              padding: "0.125rem 0.5rem",
              border: "1px solid #262626",
              fontWeight: 500,
              fontSize: "0.8125rem",
            },
            "code::before": {
              content: '""',
            },
            "code::after": {
              content: '""',
            },
            table: {
              marginTop: "3.5rem",
              marginBottom: "3.5rem",
              fontSize: "0.9375rem",
            },
            "thead th": {
              color: "#ffffff",
              fontWeight: 600,
              padding: "0.75rem 1rem",
            },
            "tbody td": {
              color: "#a3a3a3",
              padding: "0.75rem 1rem",
            },
            hr: {
              marginTop: "4.5rem",
              marginBottom: "4.5rem",
            },
            blockquote: {
              paddingLeft: "1.5rem",
              marginTop: "3rem",
              marginBottom: "3rem",
              fontStyle: "italic",
            },
            img: {
              marginTop: "3rem",
              marginBottom: "3rem",
              borderRadius: "0.5rem",
            },
            a: {
              fontWeight: 500,
              textDecoration: "none",
            },
            "a:hover": {
              color: "#6ee7b7",
            },
            strong: {
              fontWeight: 600,
            },
          },
        },
      },
    },
  },
  plugins: [typography],
};
