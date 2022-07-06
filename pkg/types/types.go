package types

import (
	"regexp"
)

// RuleType define the possible types of rules.
type RuleType string

const (
	// Set will set the value of a header.
	Set RuleType = "Set"
	// Join will concatenate the values of headers.
	Join RuleType = "Join"
	// Delete will delete the value of a header.
	Delete RuleType = "Del"
	// Rename will rename a header.
	Rename RuleType = "Rename"
	// PathEscape will net.url.PathEscape a header value to new or same header.
	PathEscape RuleType = "PathEscape"
	// PathUnEscape will net.url.PathUnEscape a header value to new or same header.
	PathUnescape RuleType = "PathUnescape"
	// QueryEscape will net.url.QueryEscape a header value to new or same header.
	QueryEscape RuleType = "QueryEscape"
	// QueryUnEscape will net.url.QueryUnEscape a header value to new or same header.
	QueryUnescape RuleType = "QueryUnescape"
	// RewriteValueRule will replace the value of a header with the provided value.
	RewriteValueRule RuleType = "RewriteValueRule"
)

// Rule struct so that we get traefik config.
type Rule struct {
	Header       string         `yaml:"Header"`       // header value
	HeaderPrefix string         `yaml:"HeaderPrefix"` // header prefix to find header
	Name         string         `yaml:"Name"`         // rule name
	Regexp       *regexp.Regexp `yaml:"-"`            // Used for rewrite, rename header matching
	Sep          string         `yaml:"Sep"`          // separator to use for join
	Type         RuleType       `yaml:"Type"`         // Differentiate rule types
	Value        string         `yaml:"Value"`
	ValueReplace string         `yaml:"ValueReplace"` // value used as replacement in rewrite
	Values       []string       `yaml:"Values"`       // values to join
}
