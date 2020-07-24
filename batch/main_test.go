package main

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestGetSuffix(t *testing.T) {
	gotenv.OverApply(strings.NewReader("ENV_NAMES=test"))
	assert.Equal(t, GetSuffixs(), []string{"test"}, "they should be equal")
}

func TestGetSuffixs(t *testing.T) {
	gotenv.OverApply(strings.NewReader("ENV_NAMES=test,all,tt"))
	assert.Equal(t, []string{"test", "all", "tt"}, GetSuffixs(), "they should be equal")

	gotenv.OverApply(strings.NewReader("ENV_NAMES=test,%all* tt"))
	assert.Equal(t, []string{"test", "all", "tt"}, GetSuffixs(), "they should be equal")

	// DO NOT USE `\` AS Delimiter
	gotenv.OverApply(strings.NewReader("ENV_NAMES=test\all3/tt"))
	assert.NotEqual(t, []string{"test", "all3", "tt"}, GetSuffixs(), "they should be equal")

	// DO NOT USE `#` AS Delimiter
	gotenv.OverApply(strings.NewReader("ENV_NAMES=test#all@tt"))
	assert.NotEqual(t, []string{"test", "all", "tt"}, GetSuffixs(), "they should be equal")
}

func TestRegexSplit(t *testing.T) {
	a := regexp.MustCompile(`[^a-zA-Z]`)
	s := a.Split("tesA%test", -1)
	t.Fatal(s)
}

func TestChk(t *testing.T) {
	n := "tes-dev"
	vg := Value{Name: &n}
	assert.Equal(t, true, Chk("dev", vg), "they should be equal")
	n = "tes-devt"
	vg = Value{Name: &n}
	assert.NotEqual(t, true, Chk("dev", vg), "they should be not equal")
}
