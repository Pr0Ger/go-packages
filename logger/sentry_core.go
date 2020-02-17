package logger

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type SentryCore struct {
	zapcore.LevelEnabler

	hub   *sentry.Hub
	scope *sentry.Scope

	BreadcrumbLevel zapcore.Level
	EventLevel      zapcore.Level
}

type SentryCoreOption func(*SentryCore)

// BreadcrumbLevel will set a minimum level of messages should be stored as breadcrumbs
func BreadcrumbLevel(level zapcore.Level) SentryCoreOption {
	return func(w *SentryCore) {
		w.BreadcrumbLevel = level
		if level > w.EventLevel {
			w.EventLevel = level
		}
	}
}

// EventLevel will set a minimum level of messages should be sent as events
func EventLevel(level zapcore.Level) SentryCoreOption {
	return func(w *SentryCore) {
		w.EventLevel = level
		if level < w.BreadcrumbLevel {
			w.BreadcrumbLevel = level
		}
	}
}

func NewSentryCore(hub *sentry.Hub, options ...SentryCoreOption) zapcore.Core {
	core := &SentryCore{
		hub:             hub,
		scope:           hub.PushScope(),
		BreadcrumbLevel: zapcore.DebugLevel,
		EventLevel:      zapcore.ErrorLevel,
	}

	for _, option := range options {
		option(core)
	}

	return core
}

func (s *SentryCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &SentryCore{
		LevelEnabler:    s.LevelEnabler,
		hub:             s.hub,
		scope:           s.hub.PushScope(),
		BreadcrumbLevel: s.BreadcrumbLevel,
		EventLevel:      s.EventLevel,
	}

	data := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(data)
	}
	clone.scope.SetExtras(data.Fields)

	return clone
}

func (s *SentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ent.Level >= s.BreadcrumbLevel {
		ce = ce.AddCore(ent, s)
	}
	return ce
}

func (s *SentryCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	data := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(data)
	}

	if ent.Level >= s.EventLevel {
		event := sentry.NewEvent()
		event.Message = ent.Message
		event.Extra = data.Fields

		s.hub.CaptureEvent(event)
	}

	breadcrumb := sentry.Breadcrumb{
		Data:      data.Fields,
		Level:     SentryLevel(ent.Level),
		Message:   ent.Message,
		Timestamp: time.Now().Unix(),
		Type:      BreadcrumbTypeDefault,
	}
	s.hub.AddBreadcrumb(&breadcrumb, nil)

	return nil
}

func (s *SentryCore) Sync() error {
	s.hub.Flush(30 * time.Second)
	return nil
}
