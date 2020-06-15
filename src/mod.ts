import { exec } from "./deps.ts";

const tem = ["tems","tes"]

console.log(tem.map(val => val + " is good").join(" + "))

await exec('kubectl version');