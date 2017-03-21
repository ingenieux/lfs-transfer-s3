package lfstransfers3

// { "event":"init", "operation":"download", "concurrent": true, "concurrenttransfers": 3 }

//GenericEvent is a generic structure to identify the underlying event
type GenericEvent struct {
	EventType string `json:"event"`
}

//InitEvent is an event to initialize the subsystem
type InitEvent struct {
	// EventType must be "init"
	EventType           string `json:"event"`
	Operation           string `json:"operation"`
	Concurrent          bool   `json:"concurrent"`
	ConcurrentTransfers int    `json:"concurrenttransfers"`
}

//TerminateEvent is sent when the system must be shutdown
type TerminateEvent struct {
	// EventType must be "terminate"
	EventType string `json:"event"`
}

//UploadEvent is sent when to upload a file
type UploadEvent struct {
	// EventType must be "upload"
	EventType string                 `json:"event"`
	Oid       string                 `json:"oid"`
	Size      uint64                 `json:"size"`
	Path      string                 `json:"path"`
	Action    map[string]interface{} `json:"action"`
}

// DownloadEvent is sent when to download a file
// { "event":"download", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245, "action": { "href": "nfs://server/path", "header": { "key": "value" } } }
type DownloadEvent struct {
	// EventType must be "download"
	EventType string                 `json:"event"`
	Oid       string                 `json:"oid"`
	Size      uint64                 `json:"size"`
	Action    map[string]interface{} `json:"action"`
}

// CompleteEvent is an outgoing event message
type CompleteEvent struct {
	// EventType must be "complete"
	EventType string      `json:"event"`
	Oid       string      `json:"oid"`
	Path      *string     `json:"path,omitempty"`
	Error     *EventError `json:"error,omitempty"`
}

// EventError contains error details
type EventError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ProgressEvent contains details about the lfs upload / download progress
// { "event":"progress", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "bytesSoFar": 1234, "bytesSinceLast": 64 }
type ProgressEvent struct {
	// EventType must be "progress"
	EventType      string `json:"event"`
	Oid            string `json:"oid"`
	BytesSoFar     uint64 `json:"bytesSoFar"`
	BytesSinceLast uint64 `json:"bytesSinceLast"`
}
