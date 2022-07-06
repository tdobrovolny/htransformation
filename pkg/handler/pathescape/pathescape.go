package pathescape

import (
	"net/http"
	"net/url"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {

	// Get all header values
	headerValues := req.Header.Values(http.CanonicalHeaderKey(rule.Value))

	req.Header.Del(rule.Header)

	for _, val := range headerValues {
		req.Header.Add(rule.Header, url.PathEscape(val))
	}
}
