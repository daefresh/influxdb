package http

import (
	"context"
	"net/http"

	"github.com/influxdata/idpe/kit/feature"
)

// Enabler allows the switching between two HTTP Handlers
type Enabler interface {
	Enabled(ctx context.Context, flagger ...feature.Flagger) bool
}
type featureHandler struct {
	enabler    Enabler
	oldHandler http.Handler
	newHandler http.Handler
}

func NewFeatureHandler(e Enabler, old, new http.Handler) *featureHandler {
	return &featureHandler{e, old, new}
}
func (c *featureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.enabler.Enabled(r.Context) {
		c.newHandler.ServeHTTP(w, r)
		return
	}
	c.oldHandler.ServeHTTP(w, r)
}
