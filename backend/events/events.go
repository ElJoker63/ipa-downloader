package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/majd/ipatool/v2/backend/models"
	"github.com/majd/ipatool/v2/backend/storage"
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

	// Device management events.
	EventDeviceConnected       EventType = "device:connected"
	EventDeviceDisconnected    EventType = "device:disconnected"
	EventDevicePairStatus      EventType = "device:pair_status"
	EventDeviceInstallProgress EventType = "device:install_progress"
	EventDeviceInstallComplete EventType = "device:install_complete"
	EventDeviceInstallFailed   EventType = "device:install_failed"
)

// Emitter defines the interface for publishing events to the frontend and persisting logs.
type Emitter interface {
	SetContext(ctx context.Context)
	SetStorage(store storage.Storage)
	Emit(eventType EventType, data ...interface{})
	EmitLog(level, message, context string)
	EmitDownloadProgress(progress models.DownloadTask)
	EmitNotification(title, message, severity string)
}

type eventEmitter struct {
	ctx     context.Context
	storage storage.Storage
	mu      sync.RWMutex
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

func (e *eventEmitter) SetStorage(store storage.Storage) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.storage = store
}

func (e *eventEmitter) Emit(eventType EventType, data ...interface{}) {
	e.mu.RLock()
	ctx := e.ctx
	e.mu.RUnlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, string(eventType), data...)
	}
}

func (e *eventEmitter) EmitLog(level, message, contextStr string) {
	e.mu.RLock()
	ctx := e.ctx
	store := e.storage
	e.mu.RUnlock()

	now := time.Now()
	entry := models.LogEntry{
		Timestamp: now,
		Level:     level,
		Message:   message,
		Context:   contextStr,
	}

	// 1. Save to SQLite database so history is preserved
	if store != nil {
		if saved, err := store.AddLog(level, message, contextStr); err == nil && saved != nil {
			entry.ID = saved.ID
		}
	}

	// 2. Console debug output
	fmt.Printf("[%s] [%s] [%s] %s\n", now.Format("15:04:05.000"), level, contextStr, message)

	// 3. Emit real-time event to Wails frontend
	if ctx != nil {
		runtime.EventsEmit(ctx, string(EventLogEntry), entry)
	}
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
