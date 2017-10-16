package model

import (
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"os"
	"time"
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
	d := new(LogField)
	return startNewLog(d)
}

func startNewLog(f *LogField) *LogField {
	return &LogField{
		TraceID: GetTraceID(),
	}
}

func initialFile(ext string) *os.File {
	f, err := os.OpenFile(getFileName(ext), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		l.Out = f
	} else {
		log.Panic("Failed to log to file, using default stderr")
	}
	l.Formatter = new(logrus.JSONFormatter)

	return f
}

func (d *LogField) Debug(m ...interface{}) {
	f := initialFile("log")
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
	f.Close()
}

func (d *LogField) Info(m ...interface{}) {
	f := initialFile("log")
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
	f.Close()
}

func (d *LogField) Panic(m ...interface{}) {
	f := initialFile("log")
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
	f.Close()
}

func (d *LogField) Warn(m ...interface{}) {
	f := initialFile("log")
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
	f.Close()
}

func (d *LogField) Log(m ...interface{}) {
	f := initialFile("log")
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
	f.Close()
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
