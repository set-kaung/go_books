<script lang="ts">
    export let id: string;
    export let name: string;
    export let link: string;
</script>

<div class="file">
    <div class="file-id">{id}</div>
    <div class="name-cont">{name}</div>
    <div><a href={link}>Download Link</a></div>
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
    a:hover {
        color: rgb(16, 129, 200);
        transform: scale(1.5);
    }
    a {
        color: rgb(57, 90, 120);
    }
    a:visited {
        color: rgb(57, 90, 120);
    }
</style>
