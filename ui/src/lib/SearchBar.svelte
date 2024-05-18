<script lang="ts">
    import { UpdateCache } from "./DataUpdate";
    import { searchTerm, duplicateMode } from "./stores";
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
</div>

<style>
    .search-bar {
        display: flex;
        gap: 0.75rem;
        font-size: 1.25rem;
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
        font-size: 0.95rem;
        padding: 0.45rem 0.75rem 0.45rem 0.75rem;
    }
    #dupes {
        font-size: 0.9rem;
    }
</style>
