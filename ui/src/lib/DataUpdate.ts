import * as DStores from "./stores";


export function UpdateCache() {
    DStores.isLoading.set(true);

    window.location.href = '/cache';
}
function wait(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function FetchFiles() {
    DStores.isLoading.set(true);
    try {
        let response = await fetch("/files");

        // If the response is a redirect, wait and try fetching again
        while (response.redirected) {
            await wait(3000); // Wait for 3 seconds before retrying
            response = await fetch("/files");
        }

        if (!response.ok) {
            if (response.status != 404) {
                DStores.isError.set(true);
            }
            throw new Error('Failed to fetch files');
        }

        const obj = await response.json();

        if ("files" in obj) {
            DStores.files.set(obj["files"]);
        } else {
            DStores.message.set(obj["message"]);
            DStores.error.set(obj["error"]);
            DStores.isError.set(true);
            console.log(obj);
        }

    } catch (err) {
        console.log(err);
        DStores.isError.set(true);
    } finally {
        DStores.isLoading.set(false); // Update isLoading when fetch completes
    }

}
