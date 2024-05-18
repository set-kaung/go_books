import File from "./File.svelte";

export function findDuplicates(files: File[]) {
    let x: { [name: string]: any; } = {};
    let dupes: any[] = [];
    files.forEach(element => {
        let name: string = element["name"];
        if (element["name"] in x) {
            dupes.push(x[name]);
            dupes.push(element);
        } else {
            x[element["name"]] = element;
        }
    });
    return dupes;
}