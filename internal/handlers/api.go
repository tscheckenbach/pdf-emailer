package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Counter holds the counter state (in-memory for demo)
type Counter struct {
	value int
	mu    sync.RWMutex
}

var counter = &Counter{value: 0}

// IncrementCounter increments the counter
func (h *Handler) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	counter.mu.Lock()
	counter.value++
	value := counter.value
	counter.mu.Unlock()

	h.logger.Info("counter incremented", "value", value)
	w.Write([]byte(strconv.Itoa(value)))
}

// DecrementCounter decrements the counter
func (h *Handler) DecrementCounter(w http.ResponseWriter, r *http.Request) {
	counter.mu.Lock()
	counter.value--
	value := counter.value
	counter.mu.Unlock()

	h.logger.Info("counter decremented", "value", value)
	w.Write([]byte(strconv.Itoa(value)))
}

// GetCounter returns the current counter value
func (h *Handler) GetCounter(w http.ResponseWriter, r *http.Request) {
	counter.mu.RLock()
	value := counter.value
	counter.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"value": value,
	})
}

// DynamicContent returns dynamic content for HTMX demo
func (h *Handler) DynamicContent(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("dynamic content requested")

	// Simulate some dynamic data
	data := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"random":    fmt.Sprintf("Random ID: %d", time.Now().Unix()),
		"status":    "Successfully loaded dynamic content!",
	}

	// Return HTML fragment for HTMX
	html := fmt.Sprintf(`
		<div class="bg-green-50 border border-green-200 rounded p-4">
			<h4 class="font-medium text-green-800 mb-2">Dynamic Content Loaded</h4>
			<div class="text-sm text-green-700 space-y-1">
				<p><strong>Timestamp:</strong> %s</p>
				<p><strong>Random ID:</strong> %s</p>
				<p><strong>Status:</strong> %s</p>
			</div>
			<button 
				hx-get="/api/dynamic-content" 
				hx-target="#dynamic-content"
				hx-swap="innerHTML"
				class="mt-3 px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700">
				Reload
			</button>
		</div>
	`, data["timestamp"], data["random"], data["status"])

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// HealthCheckJSON returns health check as JSON
func (h *Handler) HealthCheckJSON(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("health check JSON request")

	// Check database connection
	dbHealthy := true
	if err := h.db.Ping(r.Context()); err != nil {
		h.logger.Error("database health check failed", "error", err)
		dbHealthy = false
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  dbHealthy,
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
