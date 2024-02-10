import Tailwind from "tailwindcss";
import AutoPrefixer from "autoprefixer";
import Minify from "postcss-minify";

export default {
  plugins: [Tailwind, AutoPrefixer, Minify],
};
