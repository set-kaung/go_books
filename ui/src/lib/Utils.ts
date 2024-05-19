export function isExtInList(path: string, list: string[]) {
    const lastIndex = path.lastIndexOf('.');
    console.log(list);
    if (lastIndex === -1 || lastIndex === 0 || lastIndex === path.length - 1) {
        return ''; // No extension found or dot is at the beginning/end of the string
    }
    let ext = path.slice(lastIndex + 1).toLowerCase(); // Get the part after the last dot and convert to lowercase
    for (let i = 0; i < list.length; i++) {
        let item = list[i];
        console.log(ext);
        if (item == ext) {
            console.log("yes");
            return true;
        }
    }
    return false;
}