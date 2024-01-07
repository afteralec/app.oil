/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/*.html", "./web/templates/**/*.html"],
  safelist: ["text-amber-700", "text-sky-700"],
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
