import { writable } from "svelte/store";


type TimeoutId = ReturnType<typeof setInterval> | number;

export const isLoading = writable(true);
export const isError = writable(false);
export const message = writable("");
export const error = writable("");
export const duplicateMode = writable(false);
// export const elapsedTime = writable(0);
export const timer = writable<TimeoutId | null>(null);
export const searchTerm = writable("");
export const files = writable([]);
export const extensions = writable<string[]>([]);
export const extensionMode = writable(false);
export const pdfChecked = writable(false);
export const epubChecked = writable(false);
