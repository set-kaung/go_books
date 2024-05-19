<script lang="ts">
    import { UpdateCache } from "./DataUpdate";
    import {
        searchTerm,
        duplicateMode,
        extensions,
        pdfChecked,
        epubChecked,
        extensionMode,
    } from "./stores";

    let extensionMod: boolean;

    extensionMode.subscribe((v) => (extensionMod = v));

    function updateExtensionsToFilter() {
        const newExtensions = [];
        if ($pdfChecked) newExtensions.push("pdf");
        if ($epubChecked) newExtensions.push("epub");
        // Add more file extensions as needed
        extensions.set(newExtensions);
    }
</script>

<div class="search-bar">
    <div class="search">
        <label for="search">Search: </label>
        <input type="text" bind:value={$searchTerm} id="search" />
    </div>
    <button class="cache-btn" on:click={UpdateCache}>Update Cache</button>
    <div class="dupes">
        <label for="dupes">Duplicate Mode</label>
        <input
            type="checkbox"
            name="dupes"
            id="dupes"
            bind:checked={$duplicateMode}
        />
    </div>
    <div class="extensions">
        <label for="extensions">Extensions: </label>
        <input
            type="checkbox"
            name="extensions"
            id="extensions"
            bind:checked={$extensionMode}
            on:change={updateExtensionsToFilter}
        />
        {#if extensionMod}
            <div id="extension-list">
                <input
                    type="checkbox"
                    name="pdf"
                    id="pdf"
                    bind:checked={$pdfChecked}
                    on:change={updateExtensionsToFilter}
                />
                <label for="pdf">PDF</label>
                <input
                    type="checkbox"
                    name="epub"
                    id="epub"
                    bind:checked={$epubChecked}
                    on:change={updateExtensionsToFilter}
                />
                <label for="epub">Epub</label>
            </div>
        {/if}
    </div>
</div>

<style>
    :root {
        --lable-font-size: 1.15rem;
    }
    .search-bar {
        display: flex;
        gap: 0.75rem;
        font-size: var(--lable-font-size);
    }

    .search {
        display: flex;
        gap: 0.5rem;
        align-self: center;
    }
    .dupes {
        display: flex;
        gap: 0.5rem;
        align-self: center;
    }
    .cache-btn {
        padding: 0.45rem 0.75rem 0.45rem 0.75rem;
    }

    .extensions {
        display: flex;
        align-self: center;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 5px;
    }
</style>
