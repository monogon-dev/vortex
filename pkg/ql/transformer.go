package ql

import (
	"fmt"
	"log"

	"github.com/grafana/loki/pkg/logql/syntax"
	"github.com/prometheus/prometheus/model/labels"
)

const LogMessageColumn = "message"
const LogLabelsColumn = "labels"
const FilebeatTable = "filebeat"

type logQLTransformer struct {
	*selectBuilder
}

func (b *logQLTransformer) AcceptMatcher(m *labels.Matcher) {
	if m.Value == "" && m.GetRegexString() == "" {
		return
	}

	keySelector := fmt.Sprintf("arrayElement(`%s`, %s)", LogLabelsColumn, b.Args.Add(m.Name))
	switch m.Type {
	case labels.MatchEqual:
		b.Where(b.EqualField(keySelector, m.Value))
	case labels.MatchNotEqual:
		b.Where(b.NotEqualField(keySelector, m.Value))
	case labels.MatchRegexp:
		b.Where(b.MatchField(keySelector, m.GetRegexString()))
	case labels.MatchNotRegexp:
		b.Where(b.NotMatchField(keySelector, m.GetRegexString()))
	default:
		panic(fmt.Sprintf("invalid match type: %v", m.Type))
	}
}

func (b *logQLTransformer) AcceptPipelineExpr(expr *syntax.PipelineExpr) {
	for _, matcher := range expr.Matchers() {
		b.AcceptMatcher(matcher)
	}
}

func (b *logQLTransformer) AcceptMatchersExpr(expr *syntax.MatchersExpr) {
	for _, matcher := range expr.Matchers() {
		b.AcceptMatcher(matcher)
	}
}

func (b *logQLTransformer) AcceptLineFilterExpr(expr *syntax.LineFilterExpr) {
	if expr.Match == "" {
		return
	}

	switch expr.Ty {
	case labels.MatchEqual:
		b.Where(b.Like(LogMessageColumn, "%"+expr.Match+"%"))
	case labels.MatchNotEqual:
		b.Where(b.NotLike(LogMessageColumn, "%"+expr.Match+"%"))
	case labels.MatchRegexp:
		b.Where(b.Match(LogMessageColumn, expr.Match))
	case labels.MatchNotRegexp:
		b.Where(b.NotMatch(LogMessageColumn, expr.Match))
	default:
		panic(fmt.Sprintf("invalid match type: %v", expr.Ty))
	}
}

func (b *logQLTransformer) AcceptLogSelectorExpr(expr syntax.LogSelectorExpr) {
	expr.Walk(func(e any) {
		switch e.(type) {
		case *syntax.PipelineExpr:
			// ignored
		case *syntax.MatchersExpr:
			b.AcceptMatchersExpr(e.(*syntax.MatchersExpr))
		case *syntax.LineFilterExpr:
			b.AcceptLineFilterExpr(e.(*syntax.LineFilterExpr))
		default:
			log.Printf("%#v", e)
		}
	})
}
