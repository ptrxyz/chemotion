import path from "path"

const dirSort = (x, y) => {
    let [l1, l2] = [x.split(path.sep).length, y.split(path.sep).length]
    if (l1 === l2) {
        return x.localeCompare(y)
    }
    return l1 > l2 ? 1 : -1
}


export { dirSort } 