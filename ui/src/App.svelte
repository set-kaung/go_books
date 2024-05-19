<script lang="ts">
    import File from "./lib/File.svelte";
    import Loading from "./lib/Loading.svelte";
    import { findDuplicates } from "./lib/Duplicate";
    import SearchBar from "./lib/SearchBar.svelte";
    import * as DStores from "./lib/stores";
    import { onMount } from "svelte";
    import { FetchFiles } from "./lib/DataUpdate";
    import { isExtInList } from "./lib/Utils";
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
    let extensionMod;

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

    let extensionsList: string[] = [];

    DStores.extensionMode.subscribe((v) => (extensionMod = v));

    DStores.extensions.subscribe((v) => {
        extensionsList = v;
    });

    onMount(() => {
        FetchFiles();
    });

    $: {
        filesCopy = files;

        if (search_term != "") {
            filesCopy = files.filter((x: File) => {
                let fName: string = x["name"];
                return fName.toLowerCase().includes(search_term.toLowerCase());
            });
        }
        if (dupes) {
            filesCopy = findDuplicates(filesCopy);
        }
        if (extensionMod && extensionsList.length != 0) {
            filesCopy = filesCopy.filter((x: File) => {
                return isExtInList(x["name"], extensionsList);
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
