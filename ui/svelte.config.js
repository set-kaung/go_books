import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'


export default {
    // Consult https://svelte.dev/docs#compile-time-svelte-preprocess
    // for more information about preprocessors
    preprocess: vitePreprocess(),
    vitePlugin: {
        // set to true for defaults or customize with object
        inspector: {
            toggleKeyCombo: 'meta-shift',
            showToggleButton: 'always',
            toggleButtonPos: 'bottom-right'
        }
    }
}
