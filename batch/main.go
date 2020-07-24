package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/stephenafamo/kronika"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Apply(strings.NewReader("ORGS="))
	gotenv.Apply(strings.NewReader("PROJECT="))
	gotenv.Apply(strings.NewReader("API_KEY="))
	gotenv.Apply(strings.NewReader(`CRON_SCHEDULE="*/1 * * * *"`))
	gotenv.Apply(strings.NewReader("ENV_NAMES="))
}

func main() {
	version := "sycer-v0.0.1"

	fmt.Println(version)
	ctx := context.Background()

	start, err := time.Parse(
		"2006-01-02 15:04:05",
		"2019-09-17 14:00:00",
	) // any time in the past works but it should be on the hour
	if err != nil {
		panic(err)
	}
	interval := time.Second

	for t := range kronika.Every(ctx, start, interval) {
		fmt.Println(t.Format("2006-01-02 15:04:05"))
		lsCmd := exec.Command("ls", "-al")
		lsOut, err := lsCmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(lsOut))
	}
}
