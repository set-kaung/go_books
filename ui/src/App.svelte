<script lang="ts">
    import File from "./lib/File.svelte";
    import Loading from "./lib/Loading.svelte";
    import { findDuplicates } from "./lib/Duplicate";
    import SearchBar from "./lib/SearchBar.svelte";
    import * as DStores from "./lib/stores";
    import { onMount } from "svelte";
    import { FetchFiles } from "./lib/DataUpdate";

    let filesCopy;
    let search_term = "";
    let files = [];
    let isLoading = true;
    let isError = false;
    let message = "";
    let error = "";
    let dupes = false;
    let timer;
    let elapsedTime;

    DStores.duplicateMode.subscribe((dup) => {
        dupes = dup;
    });

    DStores.isError.subscribe((err) => {
        isError = err;
    });

    DStores.isLoading.subscribe((load) => {
        isLoading = load;
    });

    DStores.message.subscribe((m) => {
        message = m;
    });

    DStores.timer.subscribe((tim) => {
        timer = tim;
    });

    DStores.elapsedTime.subscribe((et) => {
        elapsedTime = et;
    });

    DStores.searchTerm.subscribe((search) => {
        search_term = search;
    });

    DStores.files.subscribe((f) => {
        files = f;
    });

    onMount(() => {
        FetchFiles();
    });

    $: {
        if (dupes) {
            filesCopy = findDuplicates(files);
        } else {
            filesCopy = files.filter((x: File) => {
                let name: string = x["name"];
                return name.toLowerCase().includes(search_term.toLowerCase());
            });
        }
    }
</script>

<div class="container">
    <SearchBar />
</div>
{#if isError}
    <div>
        <div class="message">
            <h2>{message}</h2>
        </div>
        <div class="error">
            <h3>{error}</h3>
        </div>
    </div>
{:else if isLoading}
    <Loading />
{:else}
    <div>
        <ol>
            {#each filesCopy as file (file.id)}
                <li><File {...file}></File></li>
            {/each}
        </ol>
    </div>
{/if}

<style>
    .container {
        display: flex;
        flex-direction: column;
    }
    li {
        margin-bottom: 1rem;
    }
</style>
