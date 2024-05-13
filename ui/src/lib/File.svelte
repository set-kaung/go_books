<script lang="ts">
    export let id: string;
    export let name: string;
    export let link: string;
</script>

<div class="file">
    <div class="file-id">{id}</div>
    <div class="name-cont">{name}</div>
    <div><a href={link}>Link</a></div>
</div>

<style>
    .file {
        display: flex;
    }
    .name-cont {
        margin-right: 0.3rem;
    }
    .file div {
        border: 1px solid black;
        border-collapse: collapse;
        padding: 0.5rem;
    }
    .file-id {
        display: none;
    }
</style>
