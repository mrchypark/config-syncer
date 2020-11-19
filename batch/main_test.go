package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSuffix(t *testing.T) {
	assert.Equal(t, getSuffixs("test"), []string{"test"}, "they should be equal")
}

func TestGetSuffixs(t *testing.T) {
	assert.Equal(t, []string{"test", "all", "tt"}, getSuffixs("test,all,tt"), "they should be equal")

	assert.Equal(t, []string{"test", "all", "tt"}, getSuffixs("test,%all* tt"), "they should be equal")
	assert.Equal(t, []string{"test", "all", "tt"}, getSuffixs("test#all@tt"), "they should not be equal")

	// DO NOT USE `\` AS Delimiter
	assert.NotEqual(t, []string{"test", "all3", "tt"}, getSuffixs("test\all3/tt"), "they should not be equal")
}

func TestRegexSplit(t *testing.T) {
	a := regexp.MustCompile(`[^a-zA-Z]`)
	s := a.Split("tesA%test", -1)
	t.Fatal(s)
}

func TestChk(t *testing.T) {
	n := "tes-dev"
	vg := value{Name: n}
	assert.Equal(t, true, chk("dev", vg), "they should be equal")
	n = "tes-devt"
	vg = value{Name: n}
	assert.NotEqual(t, true, chk("dev", vg), "they should be not equal")
}
