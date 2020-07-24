package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/imroc/req"
	"github.com/stephenafamo/kronika"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Apply(strings.NewReader("ORGS=SKT-AIDevOps"))
	gotenv.Apply(strings.NewReader("PROJECT=hera"))
	// azure devops PAT key
	gotenv.Apply(strings.NewReader("API_KEY=zxvxhk6a5ja7ccf4v2ehatuzejy6pwfjkfodlglw4g424zopa6ga"))
	// second
	gotenv.Apply(strings.NewReader(`INTERVAL=60`))
	gotenv.Apply(strings.NewReader("ENV_NAMES=dev"))
}

func main() {
	version := "sycer-v0.0.1"

	fmt.Println(version)
	ctx := context.Background()

	start, err := time.Parse(
		"2006-01-02 15:04:05",
		"2019-09-17 14:00:00",
	)

	if err != nil {
		panic(err)
	}

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
	envs := os.Getenv("ENV_NAMES")
	if envs == "" {
		panic("Need to set ENV_NAMES.")
	}

	key = b64.StdEncoding.EncodeToString([]byte(key))
	tar := `https://dev.azure.com/` + orgs + `/` + proj + `/_apis/distributedtask/variablegroups?api-version=5.1-preview.1`
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + key,
	}

	// b, _ := v.Marshal()
	// fmt.Println(string(b))

	s, _ := strconv.Atoi(os.Getenv("INTERVAL"))
	interval := time.Duration(s) * time.Second
	for t := range kronika.Every(ctx, start, interval) {
		fmt.Println(t.Format("2006-01-02 15:04:05"))

		res, err := req.Get(tar, authHeader)
		if err != nil {
			fmt.Println("Request fail: get variable groups")
			fmt.Println(err.Error())
			continue
		}

		v := VG{}
		res.ToJSON(&v)

		// for _, vg := range v.Value {
		// 	vg.Name
		// }

		lsCmd := exec.Command("ls", "-al")
		lsOut, err := lsCmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(lsOut))
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
	fmt.Println(s)
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
	Name      *string        `json:"name,omitempty"`
}

type Key struct {
	Value    *string `json:"value"`
	IsSecret *bool   `json:"isSecret,omitempty"`
}
