package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const versionRegexp = `Weaver v\d+.\d+.\d+ ([a-zA-Z0-9]{5,40}|N/A)\s?((\d{4})-(\d{2})-(\d{2})T(\d{2})\:(\d{2})\:(\d{2})[+-](\d{2})\:(\d{2})|SELFBUILD)?`

func TestGenVersionDefault(t *testing.T) {
	matched, err := regexp.Match(versionRegexp, []byte(GenVersion()))
	assert.NoError(t, err)
	assert.True(t, matched, "The version regex doesn't match")
}

func TestGenVersionFromCI(t *testing.T) {

	// Values supossedly injected from CI
	Build = "2018-08-09T20:11:29+07:00"
	Commit = "7800e4d69835522b712bcc7ba3242af24f6a4e4c"

	matched, err := regexp.Match(versionRegexp, []byte(GenVersion()))
	assert.NoError(t, err)
	assert.True(t, matched, "The version regex doesn't match")
}
