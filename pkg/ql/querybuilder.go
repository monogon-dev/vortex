package ql

import (
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

type selectBuilder struct {
	*sqlbuilder.SelectBuilder
}

func (b selectBuilder) Match(field, value string) string {
	return b.MatchField(sqlbuilder.Escape(field), value)
}

func (b selectBuilder) NotMatch(field, value string) string {
	return b.NotMatchField(sqlbuilder.Escape(field), value)
}

func (b selectBuilder) MatchField(field, value string) string {
	buf := &strings.Builder{}
	buf.WriteString("match(")
	buf.WriteString(field)
	buf.WriteString(", ")
	// Unlike re2's default behavior, "." matches line breaks.
	// To disable this, prepend the pattern with (?-s).
	buf.WriteString(b.Args.Add("(?-s)" + value))
	buf.WriteString(")")
	return buf.String()
}

func (b selectBuilder) NotMatchField(field, value string) string {
	return "NOT " + b.MatchField(field, value)
}

// EqualField is identical to Equal but does not escape its field
func (b selectBuilder) EqualField(field string, value interface{}) string {
	buf := &strings.Builder{}
	buf.WriteString(field)
	buf.WriteString(" = ")
	buf.WriteString(b.Args.Add(value))
	return buf.String()
}

// NotEqualField is identical to NotEqual but does not escape its field
func (b selectBuilder) NotEqualField(field string, value interface{}) string {
	buf := &strings.Builder{}
	buf.WriteString(field)
	buf.WriteString(" <> ")
	buf.WriteString(b.Args.Add(value))
	return buf.String()
}
