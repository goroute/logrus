package logrus

import (
	"strconv"
	"time"

	"github.com/goroute/route"
	"github.com/sirupsen/logrus"
)

// Options defines the options for logrus middleware.
type Options struct {
	// Skipper defines a function to skip middleware.
	Skipper route.Skipper

	// Entry defines logrus entry.
	Entry *logrus.Entry

	// Fields allows to define custom fields.
	Fields func(c route.Context, startTime, endTime time.Time) logrus.Fields
}

type Option func(*Options)

func GetDefaultOptions() Options {
	return Options{
		Skipper: route.DefaultSkipper,
		Entry:   logrus.NewEntry(logrus.StandardLogger()),
	}
}

func Skipper(skipper route.Skipper) Option {
	return func(o *Options) {
		o.Skipper = skipper
	}
}

func Entry(entry *logrus.Entry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}

func Fields(fields func(c route.Context, startTime, endTime time.Time) logrus.Fields) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

// New returns a middleware which logs requests
func New(options ...Option) route.MiddlewareFunc {
	// Apply options.
	opts := GetDefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(next route.HandlerFunc) route.HandlerFunc {
		return func(c route.Context) error {
			if opts.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// Get content length.
			contentLength := req.Header.Get(route.HeaderContentLength)
			if contentLength == "" {
				contentLength = "0"
			}

			// Create default logger fields.
			var infoFields logrus.Fields
			if opts.Fields != nil {
				infoFields = opts.Fields(c, start, stop)
			} else {
				infoFields = logrus.Fields{
					"bytes_in":   contentLength,
					"bytes_out":  strconv.FormatInt(res.Size, 10),
					"host":       req.Host,
					"start_time": start.Format(time.RFC3339),
					"end_time":   stop.Format(time.RFC3339),
					"status":     res.Status,
					"method":     req.Method,
					"path":       req.URL.Path,
					"latency_ns": strconv.FormatInt(int64(stop.Sub(start)), 10),
					"latency":    stop.Sub(start).String(),
				}
			}

			entry := opts.Entry.WithFields(infoFields)

			// Log message.
			if err != nil && res.Status >= 500 {
				entry.Error(err)
			} else {
				entry.Info()
			}

			return nil
		}
	}
}
