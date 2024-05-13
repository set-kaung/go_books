<script lang="ts">
    let search_term = "";
    import File from "./lib/File.svelte";
    let files = [];
    let isLoading = true;
    FetchFiles();
    $: filesCopy = files.filter((x: File) => {
        let name: string = x["name"];
        return name.toLowerCase().includes(search_term.toLowerCase());
    });

    function UpdateCache() {
        fetch("/cache")
            .then((resp) => {
                if (!resp.ok) {
                    throw new Error("network response was not ok");
                }
                return resp;
            })
            .catch((err) => console.log(err));
        FetchFiles();
    }

    function FetchFiles() {
        isLoading = true;
        fetch("/files")
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }

                return response.json();
            })
            .then((obj) => {
                console.log(obj);
                files = obj["files"];
            })
            .catch((err) => console.log(err))
            .finally(() => {
                isLoading = false; // Update isLoading when fetch completes
            });
    }
</script>

{#if isLoading}
    <p>Loading...</p>
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
