package server

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyMatcher(t *testing.T) {
	body := bytes.NewReader([]byte(`{ "drivers": { "id": "123", "name": "hello"} }`))

	req := httptest.NewRequest("GET", "/drivers", body)
	expr := ".drivers.id"

	key, err := matcherMux["body"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "123", key)
}

func TestBodyMatcherParseInt(t *testing.T) {
	body := bytes.NewReader([]byte(`{ "routeRequests": [{ "id": "123", "serviceType": 1}] }`))

	req := httptest.NewRequest("GET", "/drivers", body)
	expr := ".routeRequests.[0].serviceType"

	key, err := matcherMux["body"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "1", key)
}

func TestBodyMatcherParseTypeAssertFail(t *testing.T) {
	body := bytes.NewReader([]byte(`{ "routeRequests": [{ "id": "123", "serviceType": []}] }`))

	req := httptest.NewRequest("GET", "/drivers", body)
	expr := ".routeRequests.[0].serviceType"

	key, err := matcherMux["body"](req, expr)
	require.Error(t, err, "should have failed to match a key")
	require.Equal(t, "", key)

	assert.Equal(t, "failed to type assert bodyKey", err.Error())
}

func TestBodyMatcherFail(t *testing.T) {
	body := bytes.NewReader([]byte(`{ "drivers": { "id": "123", "name": "hello"} }`))

	req := httptest.NewRequest("GET", "/drivers", body)
	expr := ".drivers.blah"

	key, err := matcherMux["body"](req, expr)
	require.Error(t, err, "should have failed to match a key")

	assert.Equal(t, "", key)
}

func TestHeaderMatcher(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	req.Header.Add("Hello", "World")

	expr := "Hello"

	key, err := matcherMux["header"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "World", key)
}

func TestHeadersCsvMatcherWithSingleHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	req.Header.Add("H1", "One")
	req.Header.Add("H2", "Two")
	req.Header.Add("H3", "Three")

	expr := "H2"

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, "Two", key)
}

func TestHeadersCsvMatcherWithSingleHeaderWhenNoneArePresent(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	expr := "H1"

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, "", key)
}

func TestHeadersCsvMatcherWithZeroHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	req.Header.Add("H1", "One")
	req.Header.Add("H2", "Two")
	req.Header.Add("H3", "Three")

	expr := ""

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, "", key)
}

func TestHeadersCsvMatcherWithMultipleHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	req.Header.Add("H1", "One")
	req.Header.Add("H2", "Two")
	req.Header.Add("H3", "Three")

	expr := "H1,H3"

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, "One,Three", key)
}

func TestHeadersCsvMatcherWithMultipleHeadersWhenSomeArePresent(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	req.Header.Add("H3", "Three")

	expr := "H1,H3"

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, ",Three", key)
}

func TestHeadersCsvMatcherWithMultipleHeadersWhenNoneArePresent(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)
	expr := "H1,H3"

	key, err := matcherMux["multi-headers"](req, expr)
	require.NoError(t, err, "should not have failed to extract headers")

	assert.Equal(t, ",", key)
}

func TestHeaderMatcherFail(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers", nil)

	expr := "Hello"

	key, err := matcherMux["header"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "", key)
}

func TestParamMatcher(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers?url=blah", nil)

	expr := "url"

	key, err := matcherMux["param"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "blah", key)
}

func TestParamMatcherFail(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers?url=blah", nil)

	expr := "hello"

	key, err := matcherMux["param"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "", key)
}

func TestPathMatcher(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers/123", nil)

	expr := `/drivers/(\d+)`

	key, err := matcherMux["path"](req, expr)
	require.NoError(t, err, "should not have failed to match a key")

	assert.Equal(t, "123", key)
}

func TestPathMatcherFail(t *testing.T) {
	req := httptest.NewRequest("GET", "/drivers/123", nil)

	expr := `/drivers/blah`

	key, err := matcherMux["path"](req, expr)
	require.Error(t, err, "should have failed to match a key")

	assert.Equal(t, "", key)
}
