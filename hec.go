package hec

import (
	"io"
	"net/http"
)

type HEC interface {
	SetHTTPClient(client *http.Client)
	SetKeepAlive(enable bool)
	SetChannel(channel string)
	SetMaxRetry(retries int)
	SetMaxContentLength(size int)

	// WriteEvent writes single event via HEC json mode
	WriteEvent(event *Event) error

	// WriteBatch writes multiple events via HCE batch mode
	WriteBatch(events []*Event) error
}

// ClusteredHECRaw is an HEC interface that can submit raw event data
// to nodes in a splunk cluster.
type ClusteredHECRaw interface {
	HEC
	// WriteRaw writes raw data stream via HEC raw mode; if it
	// receives an error from the upstream server, it rewinds the
	// input reader and sends to a different server in the cluster.
	WriteRaw(reader io.ReadSeeker, metadata *EventMetadata) error
}

// HECRaw is an HEC submission interface that can submit raw event
// data.
type HECRaw interface {
	HEC
	// WriteRaw writes raw data stream via HEC raw mode
	WriteRaw(reader io.Reader, metadata *EventMetadata) error
}
