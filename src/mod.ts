import { exec } from "./deps.ts";
import ky from 'https://deno.land/x/ky/index.js';

const parsed = await ky.get('https://httpbin.org/get').json();
console.log(parsed)

const tem = ["tems","tes"]
console.log(tem.map(val => val + " is good").join(" + "))
await exec('kubectl version');