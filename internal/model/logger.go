package model

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	LogField struct {
		TraceID   string
		StartTime time.Time
		EndTime   time.Time
		Delta     float64
		Host      string
		Path      string
		Method    string
		Status    int
		Tag       string
		Service   string
	}
)

var (
	l        = logrus.New()
	Path     string
	FileName string
)

func NewLog() *LogField {
	return &LogField{
		TraceID: GetTraceID(),
	}
}

func (d *LogField) Debug(m ...interface{}) {
	l.WithFields(logrus.Fields{
		"trace-ID": d.TraceID,
		"Start":    d.StartTime,
		"End":      d.EndTime,
		"delta":    d.Delta,
		"host":     d.Host,
		"path":     d.Path,
		"service":  d.Service,
		"method":   d.Method,
		"result":   d.Status,
	}).Debug(m)
}

func (d *LogField) Info(m ...interface{}) {
	l.WithFields(logrus.Fields{
		"trace-ID": d.TraceID,
		"Start":    d.StartTime,
		"End":      d.EndTime,
		"delta":    d.Delta,
		"host":     d.Host,
		"path":     d.Path,
		"service":  d.Service,
		"method":   d.Method,
		"result":   d.Status,
	}).Info(m)
}

func (d *LogField) Panic(m ...interface{}) {
	l.WithFields(logrus.Fields{
		"trace-ID": d.TraceID,
		"Start":    d.StartTime,
		"End":      d.EndTime,
		"delta":    d.Delta,
		"host":     d.Host,
		"path":     d.Path,
		"service":  d.Service,
		"method":   d.Method,
		"result":   d.Status,
	}).Panic(m)
}

func (d *LogField) Warn(m ...interface{}) {
	l.WithFields(logrus.Fields{
		"trace-ID": d.TraceID,
		"Start":    d.StartTime,
		"End":      d.EndTime,
		"delta":    d.Delta,
		"host":     d.Host,
		"path":     d.Path,
		"service":  d.Service,
		"method":   d.Method,
		"result":   d.Status,
	}).Warn(m)
}

func (d *LogField) Log(m ...interface{}) {
	l.WithFields(logrus.Fields{
		"trace-ID": d.TraceID,
		"Start":    d.StartTime,
		"End":      d.EndTime,
		"delta":    d.Delta,
		"host":     d.Host,
		"path":     d.Path,
		"service":  d.Service,
		"method":   d.Method,
		"result":   d.Status,
	}).Info(m)
}

func GetTraceID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	chars := NUMERALS
	result := make([]byte, LN_TRACE_ID)
	for i := 0; i < LN_TRACE_ID; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func getFileName(ext string) string {
	t := time.Now()
	periode := t.Format("20060102")
	filename := FileName + "-" + periode + "." + ext
	return Path + filename
}

func (d *LogField) SetService(service string) *LogField {
	d.Service = service
	return d
}
func (d *LogField) SetTag(t string) *LogField {
	d.Tag = t
	return d
}

func (d *LogField) SetStatus(code int) *LogField {
	d.Status = code
	return d
}

func (d *LogField) SetMethod(method string) *LogField {
	d.Method = method
	return d
}
func (d *LogField) SetHost(s string) *LogField {
	d.Host = s
	return d
}

func (d *LogField) SetPath(s string) *LogField {
	d.Path = s
	return d
}

func (d *LogField) SetStart(t time.Time) *LogField {
	d.StartTime = t
	return d
}

func (d *LogField) SetEnd(t time.Time) *LogField {
	d.EndTime = t
	return d
}
func (d *LogField) SetDelta(f float64) *LogField {
	d.Delta = f
	return d
}
