import { exec, soxa, encode } from "./deps.ts";

const orgs = Deno.env.get('ORGS')
const proj = Deno.env.get('PROJECT')
const auth_key = Deno.env.get('KEY')
const pat = `:` + auth_key
const env_name = Deno.env.get('ENV_NAME')
console.log(env_name)

const vg_config = {
  baseURL: `https://dev.azure.com/${orgs}/${proj}/`,
  headers: {
    Authorization : `Basic ${encode(pat)}`
          }
 }
const parsed = await soxa.get(`_apis/distributedtask/variablegroups?api-version=5.1-preview.1`, vg_config);

console.log(parsed
  .data
  .value
  .filter((val: any) => 
    val.name.includes(`${env_name}`))
  .map((val: any) =>
    `kubectl create configmap ${val.name}${
      Object.entries(val.variables).map(([key, value]: [string, any]) => ` --from-literal=${key}=${value["value"]}`)
    } -o yaml --dry-run=client | kubectl apply -f -` ))

// await exec('kubectl version');
