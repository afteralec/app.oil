/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/*.html", "./web/templates/**/*.html"],
  safelist: [
    "light",
    "dark",
    "text-incomplete",
    "text-ready",
    "text-submitted",
    "text-review",
    "text-approved",
    "text-reviewed",
    "text-rejected",
    "text-archived",
    "text-canceled",
  ],
  theme: {
    extend: {
      colors: {
        bg: "hsl(var(--bg))",
        fg: "hsl(var(--fg))",
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          fg: "hsl(var(--primary-fg))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          fg: "hsl(var(--muted-fg))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          fg: "hsl(var(--accent-fg))",
        },
        err: {
          DEFAULT: "hsl(var(--err))",
          fg: "hsl(var(--err-fg))",
          hl: "hsl(var(--err-hl))",
        },
        warn: {
          DEFAULT: "hsl(var(--warn))",
          fg: "hsl(var(--warn-fg))",
          hl: "hsl(var(--warn-hl))",
        },
        info: {
          DEFAULT: "hsl(var(--info))",
          fg: "hsl(var(--info-fg))",
          hl: "hsl(var(--info-hl))",
        },
        success: {
          DEFAULT: "hsl(var(--success))",
          fg: "hsl(var(--success-fg))",
          hl: "hsl(var(--success-hl))",
        },
        incomplete: "hsl(var(--incomplete))",
        ready: "hsl(var(--ready))",
        submitted: "hsl(var(--submitted))",
        review: "hsl(var(--in-review))",
        reviewed: "hsl(var(--reviewed))",
        approved: "hsl(var(--approved))",
        rejected: "hsl(var(--rejected))",
        canceled: "hsl(var(--canceled))",
        archived: "hsl(var(--archived))",
      },
      maxWidth: {
        "10xl": "100rem",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
