import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import {defineConfig} from "eslint/config";


export default defineConfig([
    {
        files: ["**/*.{js,mjs,cjs,ts,mts,cts,vue}"],
        ignores: [
            "**/*.spec.*",
            "**/*.test.*"
        ],
        plugins: {js},
        extends: ["js/recommended"],
        rules: {
            "no-unused-vars": ["error", { "argsIgnorePattern": "^_" }]
        }
    },
    {
        files: ["**/*.{js,mjs,cjs,ts,mts,cts,vue}"],
        ignores: [
            "**/*.spec.*",
            "**/*.test.*"
        ],
        languageOptions: {globals: globals.browser}
    },

    pluginVue.configs["flat/essential"],
    {
        files: ["**/*.vue"],
        languageOptions: {
            parserOptions: {parser: tseslint.parser}
        },
        rules: {
            "vue/multi-word-component-names": "off",
            "@typescript-eslint/no-this-alias": "off",
            "vue/no-reserved-component-names": "off"
        }
    },
]);
