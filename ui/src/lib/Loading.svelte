<script lang="ts">
    import { elapsedTime } from "./stores";
    let elapsed: number;

    elapsedTime.subscribe((val) => {
        elapsed = val;
    });
</script>

<div class="loading">
    <h1>Loading Data</h1>
    <h3>Elaped Time: {elapsed}</h3>
</div>

<style>
    .loading {
        margin-top: 1rem;
        max-width: 30%;
        font-size: 1rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        border: solid rgb(27, 27, 27) 1px;
        box-shadow: 5px 5px 5px rgb(27, 27, 27);
    }
</style>
