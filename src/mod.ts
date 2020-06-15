import { exec, soxa } from "./deps.ts";

const parsed = await soxa.get('https://httpbin.org/get');
console.log(parsed)

const tem = ["tems","tes"]
console.log(tem.map(val => val + " is good").join(" + "))
await exec('kubectl version');