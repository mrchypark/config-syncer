package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/imroc/req"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Apply(strings.NewReader("ORGS=SKT-AIDevOps"))
	gotenv.Apply(strings.NewReader("PROJECT=hera"))
	// azure devops PAT key
	gotenv.Apply(strings.NewReader("API_KEY=nkjodwjcvtb7sdi5fd24mnso7ljkqnt3pd3dyghr33emvqbi437a"))
	gotenv.Apply(strings.NewReader("ENV_NAMES=dev"))
}

func main() {
	version := "sycer-v0.0.1"

	fmt.Println(version)

	apik := ":" + os.Getenv("API_KEY")
	if apik == ":" {
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

	c := os.Getenv("ENV_NAMES")
	envs := getSuffixs(c)
	if len(envs) == 0 {
		panic("Need to set ENV_NAMES.")
	}

	s := daprd.NewService(":8080")

	// add some input binding handler
	if err := s.AddBindingInvocationHandler("schedule", scheduleHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting service: %v", err)
	}

}

func scheduleHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("Schedule - Metadata:%v, Data:%v", in.Metadata, in.Data)

	apik := ":" + os.Getenv("API_KEY")
	orgs := os.Getenv("ORGS")
	proj := os.Getenv("PROJECT")
	c := os.Getenv("ENV_NAMES")

	apik = b64.StdEncoding.EncodeToString([]byte(apik))
	tar := `https://dev.azure.com/` + orgs + `/` + proj + `/_apis/distributedtask/variablegroups?api-version=5.1-preview.1`
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + apik,
	}

	res, err := req.Get(tar, authHeader)
	if err != nil {
		fmt.Println(err.Error())
		panic("Request fail: get variable groups")
	}

	v := vg{}
	res.ToJSON(&v)
	cmdAll := []string{}
	envs := getSuffixs(c)

	for _, vg := range v.Value {
		for _, t := range envs {
			if chk(t, vg) {
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
		cmd := exec.Command(`bash`, `-c`, v+` | kubectl apply -f -`)
		out, err := cmd.Output()
		fmt.Println("cmd")
		fmt.Println(cmd.String())
		if err != nil {
			fmt.Println("error")
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("output")
		fmt.Println(string(out))
	}

	return nil, nil
}

func getSuffixs(e string) []string {
	t := []string{}
	if e == "" {
		return t
	}
	a := regexp.MustCompile(`[^a-zA-Z0-9]`)
	s := a.Split(e, -1)
	for _, v := range s {
		if v != "" {
			t = append(t, v)
		}
	}
	return t
}

func (r *vg) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type vg struct {
	Value []value `json:"value"`
}

type value struct {
	Variables map[string]key `json:"variables,omitempty"`
	Name      string         `json:"name,omitempty"`
}

type key struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret,omitempty"`
}

func chk(envName string, vg value) bool {
	match, err := regexp.MatchString(`.`+envName+`$`, vg.Name)
	if err != nil {
		return false
	}
	return match
}
