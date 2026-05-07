package notify_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/notify"
)

// fakeSink is a controllable Sink for tests.
type fakeSink struct {
	name    string
	called  int
	errOnce bool
}

func (f *fakeSink) Name() string { return f.name }
func (f *fakeSink) Send(_ context.Context, findings []drift.Finding) error {
	f.called++
	if f.errOnce {
		return errors.New("sink error")
	}
	return nil
}

func sampleFindings() []drift.Finding {
	return []drift.Finding{
		{Container: "web", Type: drift.DriftTypeEnv, Field: "PORT", Expected: "8080", Actual: "9090"},
	}
}

func testLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestNew_NoSinks_ReturnsError(t *testing.T) {
	_, err := notify.New(log.New(bytes.NewBuffer(nil), "", 0))
	if err == nil {
		t.Fatal("expected error for zero sinks")
	}
}

func TestNew_NilLogger_ReturnsError(t *testing.T) {
	_, err := notify.New(nil, &fakeSink{name: "s"})
	if err == nil {
		t.Fatal("expected error for nil logger")
	}
}

func TestDispatch_NoFindings_SkipsSinks(t *testing.T) {
	s := &fakeSink{name: "s"}
	d, _ := notify.New(log.New(bytes.NewBuffer(nil), "", 0), s)
	if err := d.Dispatch(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.called != 0 {
		t.Fatalf("sink should not be called for empty findings, got %d calls", s.called)
	}
}

func TestDispatch_CallsAllSinks(t *testing.T) {
	s1, s2 := &fakeSink{name: "a"}, &fakeSink{name: "b"}
	d, _ := notify.New(log.New(bytes.NewBuffer(nil), "", 0), s1, s2)
	_ = d.Dispatch(context.Background(), sampleFindings())
	if s1.called != 1 || s2.called != 1 {
		t.Fatalf("expected each sink called once, got %d %d", s1.called, s2.called)
	}
}

func TestDispatch_AllSinksFail_ReturnsError(t *testing.T) {
	s := &fakeSink{name: "bad", errOnce: true}
	buf := &bytes.Buffer{}
	d, _ := notify.New(log.New(buf, "", 0), s)
	err := d.Dispatch(context.Background(), sampleFindings())
	if err == nil {
		t.Fatal("expected error when all sinks fail")
	}
}

func TestLogSink_WritesFindings(t *testing.T) {
	buf := &bytes.Buffer{}
	s := notify.NewLogSink("test", buf)
	_ = s.Send(context.Background(), sampleFindings())
	if !strings.Contains(buf.String(), "web") {
		t.Fatalf("expected container name in output, got: %s", buf.String())
	}
}

func TestSinkCount(t *testing.T) {
	s1, s2 := &fakeSink{name: "a"}, &fakeSink{name: "b"}
	d, _ := notify.New(log.New(bytes.NewBuffer(nil), "", 0), s1, s2)
	if got := d.SinkCount(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

// Ensure fakeSink satisfies notify.Sink at compile time.
var _ notify.Sink = (*fakeSink)(nil)

// silence unused import
var _ = fmt.Sprintf
