package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/savaki/jq"
)

type matcherFunc func(request *http.Request, shardExpr string) (shardKey string, err error)

var matcherMux = map[string]matcherFunc{
	"header": func(req *http.Request, expr string) (string, error) {
		return req.Header.Get(expr), nil
	},

	"multi-headers": func(req *http.Request, expr string) (string, error) {
		headers := strings.Split(expr, ",")
		var headerValues strings.Builder

		headersCount := len(headers)
		if headersCount == 0 {
			return "", nil
		}

		for idx, header := range headers {
			headerValue := req.Header.Get(header)

			headerValues.Grow(len(headerValue))
			headerValues.WriteString(headerValue)

			if (idx + 1) != headersCount {
				headerValues.Grow(1)
				headerValues.WriteString(",")
			}
		}

		return headerValues.String(), nil
	},

	"param": func(req *http.Request, expr string) (string, error) {
		return req.URL.Query().Get(expr), nil
	},

	"path": func(req *http.Request, expr string) (string, error) {
		rex := regexp.MustCompile(expr)
		match := rex.FindStringSubmatch(req.URL.Path)
		if len(match) == 0 {
			return "", fmt.Errorf("no match found for expr: %s", expr)
		}

		return match[1], nil
	},

	"body": func(req *http.Request, expr string) (string, error) {
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", errors.Wrapf(err, "failed to read request body for expr: %s", expr)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		var bodyKey interface{}
		op, err := jq.Parse(expr)
		if err != nil {
			return "", errors.Wrapf(err, "failed to parse shard expr: %s", expr)
		}

		key, err := op.Apply(requestBody)
		if err != nil {
			return "", errors.Wrapf(err, "failed to apply parsed shard expr: %s", expr)
		}

		if err := json.Unmarshal(key, &bodyKey); err != nil {
			return "", errors.Wrapf(err, "failed to unmarshal data for shard expr: %s", expr)
		}

		switch v := bodyKey.(type) {
		case string:
			return v, nil
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		default:
			return "", errors.New("failed to type assert bodyKey")
		}
	},
}
