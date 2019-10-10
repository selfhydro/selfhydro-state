package main

import (
	"net/url"
	"testing"

	"gotest.tools/assert"
)

func Test_GetSystemIDFromURLParameter(t *testing.T) {
	mockURL := &url.URL{
		RawQuery: "systemid=selfhydro",
	}
	systemID := getParameter(mockURL, "systemid")
	assert.Equal(t, "selfhydro", systemID)
}
