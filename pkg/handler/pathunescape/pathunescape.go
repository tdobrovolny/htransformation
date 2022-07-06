package pathunescape

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
		unEscaped, err := url.PathUnescape(val)
		if err != nil {
			unEscaped = err.Error()
		}
		req.Header.Add(rule.Header, unEscaped)
	}
}
