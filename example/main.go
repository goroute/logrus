package main

import (
	"errors"
	logrusMiddleware "github.com/goroute/logrus"
	"github.com/goroute/route"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	entry := logrus.NewEntry(logger)

	mux := route.NewServeMux()
	mux.Use(
		logrusMiddleware.New(
			logrusMiddleware.Entry(entry),
			logrusMiddleware.Fields(func(c route.Context, startTime, endTime time.Time) logrus.Fields {
				return logrus.Fields{
					"host":       c.Request().Host,
					"start_time": startTime.Format(time.RFC3339),
					"end_time":   endTime.Format(time.RFC3339),
				}
			}),
		),
	)

	mux.GET("/", func(c route.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	mux.GET("/err", func(c route.Context) error {
		return errors.New("ups")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
