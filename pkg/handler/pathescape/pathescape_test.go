package pathescape_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomMoulard/htransformation/pkg/handler/pathescape"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestPathEscapeHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "no transformation",
			rule: types.Rule{
				Header: "not-existing",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "transoformation without change",
			rule: types.Rule{
				Header: "X-Test",
				Value:  "Test",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Success",
			},
		},
		{
			name: "transoformation with change",
			rule: types.Rule{
				Header: "X-Test",
				Value:  "Test",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success (transform@go)",
			},
			want: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Success%20%28transform@go%29",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			pathescape.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}
