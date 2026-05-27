package events

import "time"

type BaseEvent struct {
	Occurred time.Time
}

func (e BaseEvent) OccurredAt() time.Time { return e.Occurred }

type UserCreated struct {
	BaseEvent
	UserID   string
	TenantID string
	Email    string
}

func (e UserCreated) Name() string { return "identity.user.created" }
func (e UserCreated) Version() int { return 1 }

type TenantProvisioned struct {
	BaseEvent
	TenantID   string
	TenantName string
}

func (e TenantProvisioned) Name() string { return "tenants.provisioned" }
func (e TenantProvisioned) Version() int { return 1 }

type FileUploaded struct {
	BaseEvent
	TenantID string
	ObjectID string
	Path     string
}

func (e FileUploaded) Name() string { return "storage.file.uploaded" }
func (e FileUploaded) Version() int { return 1 }

type NotificationQueued struct {
	BaseEvent
	TenantID string
	Channel  string
	Message  string
}

func (e NotificationQueued) Name() string { return "notifications.queued" }
func (e NotificationQueued) Version() int { return 1 }
