package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/loki/pkg/querier"
	"github.com/grafana/loki/pkg/util/server"
	"github.com/grafana/loki/pkg/validation"
	"github.com/weaveworks/common/user"

	chquerier "github.com/monogon-dev/vortex/pkg/querier"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Protocol: clickhouse.HTTP,
		Addr:     []string{"127.0.0.1:8123"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	defaultLimits := validation.Limits{}
	flagext.DefaultValues(&defaultLimits)
	limits, err := validation.NewOverrides(defaultLimits, nil)
	if err != nil {
		log.Fatal(err)
	}

	chq := chquerier.NewClickhouseQuerier(conn)
	api := querier.NewQuerierAPI(querier.Config{}, chq, limits, logger)

	routes := map[string]http.Handler{
		"/loki/api/v1/query_range": querier.WrapQuerySpanAndTimeout("query.RangeQuery", api).Wrap(http.HandlerFunc(api.RangeQueryHandler)),
		"/loki/api/v1/query":       querier.WrapQuerySpanAndTimeout("query.InstantQuery", api).Wrap(http.HandlerFunc(api.InstantQueryHandler)),

		"/loki/api/v1/label":               http.HandlerFunc(api.LabelHandler),
		"/loki/api/v1/labels":              http.HandlerFunc(api.LabelHandler),
		"/loki/api/v1/label/{name}/values": http.HandlerFunc(api.LabelHandler),

		"/loki/api/v1/series":              querier.WrapQuerySpanAndTimeout("query.Series", api).Wrap(http.HandlerFunc(api.SeriesHandler)),
		"/loki/api/v1/index/stats":         querier.WrapQuerySpanAndTimeout("query.IndexStats", api).Wrap(http.HandlerFunc(api.IndexStatsHandler)),
		"/loki/api/v1/index/series_volume": querier.WrapQuerySpanAndTimeout("query.SeriesVolume", api).Wrap(http.HandlerFunc(api.SeriesVolumeHandler)),
	}

	router := mux.NewRouter()
	for path, handler := range routes {
		path, handler := path, handler
		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			newCtx := user.InjectOrgID(r.Context(), "0")
			if err := r.ParseForm(); err != nil {
				server.WriteError(err, w)
				return
			}
			handler.ServeHTTP(w, r.WithContext(newCtx))
		})
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		http.NotFoundHandler().ServeHTTP(w, r)
	})

	log.Println("Started")
	log.Fatal(http.ListenAndServe(":3100", router))
}
