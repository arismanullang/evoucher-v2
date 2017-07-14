package model

import (
	"os"
	"log"
	"github.com/sirupsen/logrus"
	"time"
	"math/rand"
)
type(
	LogField struct {
		TraceID	string
		Time	time.Time
		Delta	float64
		Service	string
		Method	string
		Tag	string
		Status	int
	}
)

var (
	l = logrus.New()
	Path string
	FileName string

)

func NewLog() *LogField{
	d := LogField{}
	d.TraceID = getTraceID()
	d.Time = time.Now()

	return startNewLog(d)
}


func startNewLog(f LogField) *LogField{
	return &LogField{
		TraceID: f.TraceID,
		Time: f.Time,
	}
}


func initialFile(ext string) *os.File{
	f, err := os.OpenFile(getFileName(ext), os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		l.Out = f
	} else {
		log.Panic("Failed to log to file, using default stderr")
	}
	l.Formatter = new(logrus.JSONFormatter)

	return f
}

func (d *LogField) Debug(m ...interface{}){
	f := initialFile("log")
	l.WithFields(logrus.Fields{
		"trace-ID" : d.TraceID,
		"time": d.Time,
		"delta": d.getDeltaTime(),
		"service": d.Service,
		"tag": d.Tag,
		"method": d.Method,
		"result":d.Status,
	}).Debug(m)
	f.Close()
}

func (d *LogField) Info(m ...interface{}){
	f := initialFile("log")
	l.WithFields(logrus.Fields{
		"trace-ID" : d.TraceID,
		"time": d.Time,
		"delta": d.getDeltaTime(),
		"service": d.Service,
		"tag": d.Tag,
		"method": d.Method,
		"result":d.Status,
	}).Info(m)
	f.Close()
}

func (d *LogField) Panic(m ...interface{}){
	f := initialFile("log")
	l.WithFields(logrus.Fields{
		"trace-ID" : d.TraceID,
		"time": d.Time,
		"delta": d.getDeltaTime(),
		"service": d.Service,
		"tag": d.Tag,
		"method": d.Method,
		"result":d.Status,
	}).Panic(m)
	f.Close()
}

func (d *LogField) Warn(m ...interface{}){
	f := initialFile("log")
	l.WithFields(logrus.Fields{
		"trace-ID" : d.TraceID,
		"time": d.Time,
		"delta": d.getDeltaTime(),
		"service": d.Service,
		"tag": d.Tag,
		"method": d.Method,
		"result":d.Status,
	}).Warn(m)
	f.Close()
}

func (d *LogField) Log(m ...interface{}){
	f := initialFile("log")
	l.WithFields(logrus.Fields{
		"trace-ID" : d.TraceID,
		"time": d.Time,
		"delta": d.getDeltaTime(),
		"service": d.Service,
		"tag": d.Tag,
		"method": d.Method,
		"result":d.Status,
	}).Info(m)
	f.Close()
}


func getTraceID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	chars := NUMERALS
	result := make([]byte, LN_TRACE_ID)
	for i := 0; i < LN_TRACE_ID; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func getFileName(ext string) string{
	t := time.Now()
	periode := t.Format("20060102")
	filename :=FileName+"-"+periode+"."+ext
	return Path+filename
}

func (d *LogField) getDeltaTime() float64 {
	return time.Since(d.Time).Seconds()
}

func (d *LogField) SetService(service string) *LogField {
	d.Service = service
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

func (d *LogField) SetTag(tag string) *LogField {
	d.Tag = tag
	return d
}
