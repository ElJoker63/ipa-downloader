package events

import (
	"context"
	"sync"
	"time"

	"github.com/majd/ipatool/v2/backend/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// EventType defines system-wide event names.
type EventType string

const (
	EventAuthStatus        EventType = "auth:status"
	EventAuthAccount       EventType = "auth:account"
	EventDownloadProgress  EventType = "download:progress"
	EventDownloadStatus    EventType = "download:status"
	EventDownloadCompleted EventType = "download:completed"
	EventDownloadFailed    EventType = "download:failed"
	EventDownloadCancelled EventType = "download:cancelled"
	EventDownloadPaused    EventType = "download:paused"
	EventLogEntry          EventType = "log:entry"
	EventNotification      EventType = "notification:show"
	EventFavoritesUpdated  EventType = "favorites:updated"
)

// Emitter defines the interface for publishing events to the frontend.
type Emitter interface {
	SetContext(ctx context.Context)
	Emit(eventType EventType, data ...interface{})
	EmitLog(level, message, context string)
	EmitDownloadProgress(progress models.DownloadTask)
	EmitNotification(title, message, severity string)
}

type eventEmitter struct {
	ctx context.Context
	mu  sync.RWMutex
}

// NewEmitter creates a new event dispatcher.
func NewEmitter() Emitter {
	return &eventEmitter{}
}

func (e *eventEmitter) SetContext(ctx context.Context) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ctx = ctx
}

func (e *eventEmitter) Emit(eventType EventType, data ...interface{}) {
	e.mu.RLock()
	ctx := e.ctx
	e.mu.RUnlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, string(eventType), data...)
	}
}

func (e *eventEmitter) EmitLog(level, message, context string) {
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Context:   context,
	}
	e.Emit(EventLogEntry, entry)
}

func (e *eventEmitter) EmitDownloadProgress(task models.DownloadTask) {
	e.Emit(EventDownloadProgress, task)
}

func (e *eventEmitter) EmitNotification(title, message, severity string) {
	payload := map[string]string{
		"title":    title,
		"message":  message,
		"severity": severity, // "info", "success", "warning", "error"
	}
	e.Emit(EventNotification, payload)
}
