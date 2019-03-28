package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gojektech/weaver/pkg/logger"
)

func Render(o interface{}) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(o)
	if err != nil {
		logger.Fatalf("Error marshaling outptu: %s", err)
		panic(err)
	}
	fmt.Print(string(buffer.Bytes()))
}
