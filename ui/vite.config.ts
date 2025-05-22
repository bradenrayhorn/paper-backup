import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { viteSingleFile } from "vite-plugin-singlefile";
import license from "rollup-plugin-license";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
    viteSingleFile(),
    {
      ...license({
        sourcemap: true,
        thirdParty: {
          multipleVersions: true,
          allow: {
            test: "(MIT OR ISC)",
            failOnViolation: true,
            failOnUnlicensed: true,
          },
          output: {
            file: path.join(__dirname, "src", "license", "npm.txt"),
          },
        },
      }),
      enforce: "post",
      apply: () => {
        return !!process.env.GENERATE_LICENSES;
      },
    },
  ],
  esbuild: {
    supported: {
      "top-level-await": true,
    },
  },
});
