// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file
package xopjson

import (
	"io"
	"sync/atomic"
	"time"

	"github.com/muir/xoplog/trace"
	"github.com/muir/xoplog/xopbase"
	"github.com/muir/xoplog/xopconst"
	"github.com/muir/xoplog/xoputil"
)

const maxBufferToKeep = 1024 * 10

var (
	_ xopbase.Logger  = &Logger{}
	_ xopbase.Request = &Span{}
	_ xopbase.Span    = &Span{}
	_ xopbase.Line    = &Line{}
)

type Option func(*Logger)

type timeOption int

const (
	noTime timeOption = iota
	strftimeTime
	timeTime
	epochTime
)

type Logger struct {
	w             io.Writer
	timeOption    timeOption
	timeFormat    string
	callers       int
	withGoroutine bool
	flusher       func()
}

type Request struct {
	*Span
}

type prefill struct {
	data string
	msg  string
}

type Span struct {
	attributes xoputil.AttributeBuilder
	trace      trace.Bundle
	logger     *Logger
	prefill    atomic.Value
}

type Line struct {
	dataBuffer xoputil.JBuilder
	level      xopconst.Level
	timestamp  time.Time
	span       *Span
}

// WithStrftime adds a timestamp to each log line.  See
// https://github.com/phuslu/fasttime for the supported
// formats.
func WithStrftime(format string) Option {
	return func(l *Logger) {
		l.timeOption = strftimeTime
		l.timeFormat = format
	}
}

func WithCallersAtLevel(logLevel xopconst.Level, framesWanted int) Option {
	return func(l *Logger) {
		l.callers = levels
	}
}

// WithFlusher should be used if the io.Writer is buffering it's output
// and can be flushed.
func WithFlusher(flusher func() Option) {
	return func(l *Logger) {
		l.flusher = flusher
	}
}

func WithGoroutineID(b bool) Option {
	return func(l *Logger) {
		l.withGoroutine = b
	}
}

func New(w io.Writer, opts ...Option) *Logger {
	logger := &Logger{
		writer: w,
	}
	for _, f := range opts {
		f(logger)
	}
	return logger
}

func (l *Logger) Close()         {}
func (l *Logger) Buffered() bool { return l.flusher != nil }

func (l *Logger) ReferencesKept() bool { return false }

func (l *Logger) Request(span trace.Bundle, name string) xopbase.Request {
	s := &Span{
		logger: l,
	}
	s.Attributes.Reset()
	return s
}

func (s *Span) Flush() {
	if s.logger.flusher != nil {
		s.logger.flusher()
	}
}

func (s *Span) Boring(bool) {} // TODO

func (s *Span) Span(span trace.Bundle, name string) xopbase.Span {
	return l.Request(span, name)
}

func (s *Span) getPrefill() *prefill {
	p := s.prefill.Load()
	if p == nil {
		return nil
	}
	return p.(*prefill)
}

func (s *Span) Line(level xopconst.Level, t time.Time) xopbase.Line {
	l := &Line{
		level:     level,
		timestamp: t,
		span:      s,
	}
	l.getPrefill()
	return l
}

func (l *Line) Recycle(level xopconst.Level, t time.Time) {
	l.level = level
	l.timestamp = t
	l.dataBuffer.Reset()
	l.getPrefill()
	return l
}

func (l *Line) getPrefill() {
	l.prefill = s.GetPrefill()
	_, _ = l.dataBuffer.WriteByte('{') // }
	if l.prefill != nil {
		l.dataBuffer.Write(l.prefill.data)
	}
}

func (l *Line) SetAsPrefill(m string) {
	skip := 1
	if l.prefill != nil {
		// don't include the prefill for _this_ line in the new prefill
		skip += len(l.prefill.data)
	}
	prefill := prefill{
		msg:  m,
		data: l.dataBuffer.String()[skip:],
	}
	l.Span.prefill.Store(prefill)
	// this Line will not be recycled so destory its buffers
	l.reclaimMemory()
}

func (l *Line) Msg(m string) {
	l.dataBuffer.Comma()
	_, _ = l.dataBuffer.WriteString(`"msg":`)
	l.dataBuffer.String(m)
	// {
	_, _ = l.dataBuffer.WriteByte("}")
	_, _ = l.span.logger.Write(l.dataBuffer.Bytes())
	l.reclaimMemory()
}

func (l *Line) reclaimMemory() {
	if l.databuffer.Len() > maxBufferToKeep {
		l.dataBuffer = xoputil.JBuffer{}
	}
}

func (l *Line) Template(m string) {
	l.dataBuffer.Comma()
	_, _ = l.dataBuffer.WriteString(`"xop":"template","msg":`)
	l.dataBuffer.String(m)
	// {
	_, _ = l.dataBuffer.WriteByte("}")
	_, _ = l.span.logger.Write(l.dataBuffer.Bytes())
	l.reclaimMemory()
}

func (l *Line) Any(k string, v interface{}) {
	// XXX
}

func (l *Line) Enum(k *xopconst.EnumAttribute, v xopconst.Enum) {
	// XXX
}

func (l *Line) Time(k string, v time.Time) { // XXX
}

func (l *Line) Link(k string, v trace.Trace) { // XXX
}

func (l *Line) Bool(k string, v bool) {
	l.dataBuffer.Comma()
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"')
	l.dataBuffer.Buf = append(l.databuffer.Buf, k...)
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"', ':')
	l.dataBuffer.Bool(v)
}

func (l *Line) Int64(k string, v int64) {
	l.dataBuffer.Comma()
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"')
	l.dataBuffer.Buf = append(l.databuffer.Buf, k...)
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"', ':')
	l.dataBuffer.Int64(v)
}

func (l *Line) Str(k string, v string) {
	l.dataBuffer.Comma()
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"')
	l.dataBuffer.Buf = append(l.databuffer.Buf, k...)
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"', ':')
	l.dataBuffer.String(v)
}

func (l *Line) Number(k string, v string) {
	l.dataBuffer.Comma()
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"')
	l.dataBuffer.Buf = append(l.databuffer.Buf, k...)
	l.dataBuffer.Buf = append(l.databuffer.Buf, '"', ':')
	l.dataBuffer.Float64(v)
}

func (s *Span) MetadataAny(k *xopconst.AnyAttribute, v interface{}) { s.Attributes.MetadataAny(k, v) }
func (s *Span) MetadataBool(k *xopconst.BoolAttribute, v bool)      { s.Attributes.MetadataBool(k, v) }
func (s *Span) MetadataEnum(k *xopconst.EnumAttribute, v xopconst.Enum) {
	s.Attributes.MetadataEnum(k, v)
}
func (s *Span) MetadataInt64(k *xopconst.Int64Attribute, v int64) { s.Attributes.MetadataInt64(k, v) }
func (s *Span) MetadataLink(k *xopconst.LinkAttribute, v trace.Trace) {
	s.Attributes.MetadataLink(k, v)
}
func (s *Span) MetadataNumber(k *xopconst.NumberAttribute, v float64) {
	s.Attributes.MetadataNumber(k, v)
}
func (s *Span) MetadataStr(k *xopconst.StrAttribute, v string)      { s.Attributes.MetadataStr(k, v) }
func (s *Span) MetadataTime(k *xopconst.TimeAttribute, v time.Time) { s.Attributes.MetadataTime(k, v) }

// end
