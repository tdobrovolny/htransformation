package htransformation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/tomMoulard/htransformation/pkg/handler/deleter"
	"github.com/tomMoulard/htransformation/pkg/handler/join"
	"github.com/tomMoulard/htransformation/pkg/handler/pathescape"
	"github.com/tomMoulard/htransformation/pkg/handler/pathunescape"
	"github.com/tomMoulard/htransformation/pkg/handler/queryescape"
	"github.com/tomMoulard/htransformation/pkg/handler/queryunescape"
	"github.com/tomMoulard/htransformation/pkg/handler/rename"
	"github.com/tomMoulard/htransformation/pkg/handler/rewrite"
	"github.com/tomMoulard/htransformation/pkg/handler/set"
	"github.com/tomMoulard/htransformation/pkg/types"
)

// HeadersTransformation holds the necessary components of a Traefik plugin.
type HeadersTransformation struct {
	name  string
	next  http.Handler
	rules []types.Rule
}

// Config holds configuration to be passed to the plugin.
type Config struct {
	Rules []types.Rule
}

// CreateConfig populates the Config data object.
func CreateConfig() *Config {
	return &Config{
		Rules: []types.Rule{},
	}
}

var ruleHandlers = map[types.RuleType]func(http.ResponseWriter, *http.Request, types.Rule){
	types.Delete:           deleter.Handle,
	types.Join:             join.Handle,
	types.PathEscape:       pathescape.Handle,
	types.PathUnescape:     pathunescape.Handle,
	types.QueryEscape:      queryescape.Handle,
	types.QueryUnescape:    queryunescape.Handle,
	types.Rename:           rename.Handle,
	types.RewriteValueRule: rewrite.Handle,
	types.Set:              set.Handle,
}

var errMissingRequiredFields = errors.New("missing required fields")

var errInvalidRuleType = errors.New("invalid rule type")

var errInvalidRegexp = errors.New("invalid regexp")

// New instantiates and returns the required components used to handle an HTTP request.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	for i, rule := range config.Rules {
		if _, ok := ruleHandlers[rule.Type]; !ok {
			return nil, fmt.Errorf("%w: %s", errInvalidRuleType, rule.Name)
		}

		if rule.Type == types.Join && (len(rule.Values) == 0 || rule.Sep == "") {
			return nil, fmt.Errorf("%w for rule %q", errMissingRequiredFields, rule.Name)
		}

		if rule.Type == types.RewriteValueRule && rule.ValueReplace == "" {
			return nil, fmt.Errorf("%w for rule %q", errMissingRequiredFields, rule.Name)
		}

		if rule.Type == types.Rename || rule.Type == types.RewriteValueRule {
			var err error
			config.Rules[i].Regexp, err = regexp.Compile(rule.Header)

			if err != nil {
				return nil, fmt.Errorf("%w: %s", errInvalidRegexp, rule.Name)
			}
		}
	}

	return &HeadersTransformation{
		name:  name,
		next:  next,
		rules: config.Rules,
	}, nil
}

// Iterate over every header to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	for _, rule := range u.rules {
		h, ok := ruleHandlers[rule.Type]
		if !ok {
			continue
		}

		h(responseWriter, request, rule)
	}

	u.next.ServeHTTP(responseWriter, request)
}
