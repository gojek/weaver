package views_test

import (
	"github.com/gojektech/weaver/internal/views"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestShouldPrettyPrintAnyJsonEnabledStruct(t *testing.T) {
	realStdout := os.Stdout
	reader, fakeStdout, err := os.Pipe()
	assert.NoError(t, err, "Error in setting fake stdout")

	os.Stdout = fakeStdout
	defer func() { os.Stdout = realStdout }()
	views.Render(struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{"gowtham", 23})

	fakeStdoutCloseErr := fakeStdout.Close()
	assert.NoError(t, fakeStdoutCloseErr, "Error close fake stdout")

	outputBuffer, err := ioutil.ReadAll(reader)
	assert.NoError(t, err, "Error in reading output from fake stdout")
	assert.Equal(t, string(outputBuffer),
		`{
    "name": "gowtham",
    "age": 23
}
`)
}
