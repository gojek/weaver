package views

import (
	"encoding/json"
	"fmt"
	"github.com/gojektech/weaver/pkg/logger"
)

func Render(o interface{}) {
	jsonBytes, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		logger.Fatalf("Error marshaling outptu: %s", err)
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
