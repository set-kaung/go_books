<script lang="ts">
    let search_term = "";
    import File from "./lib/File.svelte";
    let files = [];
    let isLoading = true;
    let isError = false;
    let message = "";
    let error = "";
    FetchFiles();
    $: filesCopy = files.filter((x: File) => {
        let name: string = x["name"];
        return name.toLowerCase().includes(search_term.toLowerCase());
    });

    function UpdateCache() {
        fetch("/cache")
            .then((resp) => {
                if (!resp.ok) {
                    if (resp.status != 404) {
                        isError = true;
                    }

                    throw new Error("network response was not ok");
                }
                return resp.json();
            })
            .then((obj) => {
                if ("error" in obj) {
                    message = obj["message"];
                    error = obj["error"];
                    isError = true;
                } else {
                    isError = false;
                }
            });
        FetchFiles();
    }

    function FetchFiles() {
        isLoading = true;
        fetch("/files")
            .then((response) => {
                if (!response.ok) {
                    if (response.status != 404) {
                        isError = true;
                    }
                }
                return response.json();
            })
            .then((obj) => {
                if ("files" in obj) {
                    files = obj["files"];
                    isError = false;
                } else {
                    message = obj["message"];
                    error = obj["error"];
                    console.log(obj);
                }
            })
            .catch((err) => {
                console.log(err);
            })
            .finally(() => {
                isLoading = false; // Update isLoading when fetch completes
            });
    }
</script>

{#if isLoading}
    <p>Loading...</p>
{:else if isError}
    <div>
        <div class="message">
            <h2>{message}</h2>
        </div>
        <div class="error">
            <h3>{error}</h3>
        </div>
    </div>
{:else}
    <div class="container">
        <div class="search-bar">
            <label for="search">Search: </label>
            <input type="text" bind:value={search_term} id="search" />
            <button on:click={UpdateCache}>Update Cache</button>
        </div>
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
    .search-bar {
        font-size: 1.25rem;
    }
</style>
