package ql

import (
	"fmt"
	"time"

	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql/syntax"
	"github.com/huandu/go-sqlbuilder"
	"github.com/prometheus/prometheus/model/labels"
)

func SelectLogsQuery(expr syntax.LogSelectorExpr, start time.Time, end time.Time, limit uint32, direction logproto.Direction) (string, []any) {
	sb := &selectBuilder{sqlbuilder.ClickHouse.NewSelectBuilder()}
	sb.Select("*").From(FilebeatTable)

	sb.Where(fmt.Sprintf("`@timestamp` >= fromUnixTimestamp64Milli(%d)", start.UnixMilli()))
	sb.Where(fmt.Sprintf("`@timestamp` <= fromUnixTimestamp64Milli(%d)", end.UnixMilli()))
	sb.Limit(int(limit))

	orderBy := "`@timestamp`"
	switch direction {
	case logproto.BACKWARD:
		orderBy += " DESC"
	case logproto.FORWARD:
		orderBy += " ASC"
	}
	sb.OrderBy(orderBy)

	l := &logQLTransformer{sb}
	l.AcceptLogSelectorExpr(expr)

	return sb.Build()
}

func LabelQuery(name string, values bool, start *time.Time, end *time.Time) (string, []any) {
	sb := &selectBuilder{sqlbuilder.ClickHouse.NewSelectBuilder()}
	sb.From(FilebeatTable).Distinct()

	if values {
		sb.Select(fmt.Sprintf("arrayElement(`labels`, %s)", sb.Args.Add(name)))
	} else {
		sb.Select("arrayJoin(mapKeys(`labels`))")
	}
	if start != nil {
		sb.Where(fmt.Sprintf("`@timestamp` >= fromUnixTimestamp64Milli(%d)", start.UnixMilli()))
	}
	if end != nil {
		sb.Where(fmt.Sprintf("`@timestamp` <= fromUnixTimestamp64Milli(%d)", end.UnixMilli()))
	}

	return sb.Build()
}

func SeriesQuery(groups [][]*labels.Matcher, start time.Time, end time.Time) (string, []any) {
	sb := &selectBuilder{sqlbuilder.ClickHouse.NewSelectBuilder()}
	sb.From(FilebeatTable).Distinct()

	sb.Select("`labels`")
	sb.Where(fmt.Sprintf("`@timestamp` >= fromUnixTimestamp64Milli(%d)", start.UnixMilli()))
	sb.Where(fmt.Sprintf("`@timestamp` <= fromUnixTimestamp64Milli(%d)", end.UnixMilli()))

	l := &logQLTransformer{sb}

	for _, group := range groups {
		for _, matcher := range group {
			l.AcceptMatcher(matcher)
		}
	}

	return l.Build()
}
