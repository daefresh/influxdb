package http

import (
	"context"
	"net/http"

	"github.com/influxdata/influxdb/v2/kit/feature"
)

// needs to be a "resource handler"

// Enabler allows the switching between two HTTP Handlers
type Enabler interface {
	Enabled(ctx context.Context, flagger ...feature.Flagger) bool
}
type featureHandler struct {
	enabler    Enabler
	oldHandler http.Handler
	newHandler http.Handler
	prefix     string
}

func NewFeatureHandler(e Enabler, old, new http.Handler, prefix string) *featureHandler {
	return &featureHandler{e, old, new, prefix}
}

func (h *featureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.enabler.Enabled(r.Context()) {
		h.newHandler.ServeHTTP(w, r)
		return
	}
	h.oldHandler.ServeHTTP(w, r)
}

func (h *featureHandler) Prefix() string {
	return h.prefix
}
