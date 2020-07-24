package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/imroc/req"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Apply(strings.NewReader("ORGS=SKT-AIDevOps"))
	gotenv.Apply(strings.NewReader("PROJECT=hera"))
	// azure devops PAT key
	gotenv.Apply(strings.NewReader("API_KEY="))
	gotenv.Apply(strings.NewReader("ENV_NAMES=dev"))
}

func main() {
	version := "sycer-v0.0.1"

	fmt.Println(version)

	key := ":" + os.Getenv("API_KEY")
	if key == ":" {
		panic("Need to set API_KEY.")
	}
	orgs := os.Getenv("ORGS")
	if orgs == "" {
		panic("Need to set ORGS.")
	}
	proj := os.Getenv("PROJECT")
	if proj == "" {
		panic("Need to set PROJECT.")
	}
	envs := GetSuffixs()
	if len(envs) == 0 {
		panic("Need to set ENV_NAMES.")
	}

	key = b64.StdEncoding.EncodeToString([]byte(key))
	tar := `https://dev.azure.com/` + orgs + `/` + proj + `/_apis/distributedtask/variablegroups?api-version=5.1-preview.1`
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + key,
	}

	res, err := req.Get(tar, authHeader)
	if err != nil {
		fmt.Println(err.Error())
		panic("Request fail: get variable groups")
	}

	v := VG{}
	res.ToJSON(&v)
	cmdAll := []string{}

	for _, vg := range v.Value {
		for _, t := range envs {
			if Chk(t, vg) {
				re := regexp.MustCompile(`.` + t + `$`)
				vn := re.ReplaceAllString(vg.Name, "")
				cmd := "kubectl create configmap " + vn
				for k, v := range vg.Variables {
					cmd += ` --from-literal=` + k + `=` + v.Value
				}
				cmd += ` -o yaml --dry-run=client`
				cmdAll = append(cmdAll, cmd)
			}
		}
	}
	for _, v := range cmdAll {
		fmt.Println(`/bin/bash` + ` -c "` + v + ` | kubectl apply -f -"`)
		cmd := exec.Command(`/bin/bash`, `-c "`+v+` | kubectl apply -f -"`)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(string(out))
	}
}

func GetSuffixs() []string {
	c := os.Getenv("ENV_NAMES")
	t := []string{}
	if c == "" {
		return t
	}
	a := regexp.MustCompile(`[^a-zA-Z0-9]`)
	s := a.Split(c, -1)
	for _, v := range s {
		if v != "" {
			t = append(t, v)
		}
	}
	return t
}

func (r *VG) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VG struct {
	Value []Value `json:"value"`
}

type Value struct {
	Variables map[string]Key `json:"variables,omitempty"`
	Name      string         `json:"name,omitempty"`
}

type Key struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret,omitempty"`
}

func Chk(envName string, vg Value) bool {
	match, err := regexp.MatchString(`.`+envName+`$`, vg.Name)
	if err != nil {
		return false
	}
	return match
}
