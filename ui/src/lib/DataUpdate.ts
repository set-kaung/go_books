import * as DStores from "./stores";

function startTimer() {
    DStores.timer.set(
        setInterval(() => {
            DStores.elapsedTime.update(n => n + 1);
        }, 1000) as unknown as number
    );
}

function stopTimer() {
    // Clear the interval and reset the timer store
    DStores.timer.update(id => {
        if (id !== null) {
            clearInterval(id as number);
            DStores.elapsedTime.set(0);
        }
        return null;
    });
}

export function UpdateCache() {

    DStores.isLoading.set(true);
    startTimer();

    fetch("/cache")
        .then((resp) => {
            if (!resp.ok) {
                if (resp.status != 404) {
                    DStores.isError.set(true);
                }

                throw new Error("network response was not ok");
            }
            return resp.json();
        })
        .then((obj) => {
            if ("error" in obj) {
                DStores.message.set(obj["message"]);
                DStores.error.set(obj["error"]);
                DStores.isError.set(true);
            } else {
                DStores.isError.set(false);
            }
        })
        .catch(err => console.log(err))
        .finally(() => FetchFiles());
}

export function FetchFiles() {
    DStores.isLoading.set(true);
    fetch("/files")
        .then((response) => {
            if (!response.ok) {
                if (response.status != 404) {
                    DStores.isError.set(true);
                }
            }
            return response.json();
        })
        .then((obj) => {
            if ("files" in obj) {
                DStores.files.set(obj["files"]);
                DStores.isError.set
            } else {
                DStores.message.set(obj["message"]);
                DStores.error.set(obj["error"]);
                DStores.isError.set(true);
                console.log(obj);
            }
        })
        .catch((err) => {
            console.log(err);
        })
        .finally(() => {
            DStores.isLoading.set(false); // Update isLoading when fetch completes
            stopTimer();
        });
}
