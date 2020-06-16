import { exec, soxa, encode } from "./deps.ts";

const orgs = Deno.env.get("ORGS");
const proj = Deno.env.get("PROJECT");
const auth_key = Deno.env.get("KEY");
const pat = `:` + auth_key;
const env_name = Deno.env.get("ENV_NAME");

// await exec('kubectl version');

const vg_config = {
  baseURL: `https://dev.azure.com/${orgs}/${proj}/`,
  headers: {
    Authorization: `Basic ${encode(pat)}`,
  },
};

const parsed = await soxa.get(
  `_apis/distributedtask/variablegroups?api-version=5.1-preview.1`,
  vg_config
);

parsed.data.value
  .filter((val: any) => val.name.includes(`drone`))
  .map(
    (val: any) =>
      `kubectl create configmap ${val.name.replace(
        `-${env_name}`,
        ""
      )}${Object.entries(val.variables)
        .map(
          ([key, value]: [string, any]) =>
            ` --from-literal=${key}=${value["value"]}`
        )
        .join("")} -o yaml --dry-run=client`
  )
  .map(async (val: string) => {
    console.log(await exec(`bash -c "${val} | kubectl apply -f -"`));
  });
