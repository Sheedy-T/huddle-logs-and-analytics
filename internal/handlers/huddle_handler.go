package handlers

import (
	"encoding/json"
	"net/http"
	"strings" // <-- This is the package required for the split helper function

	"github.com/Sheedy-T/huddle-backend/internal/core/ports"
)

type HuddleHandler struct {
	service ports.HuddleService
}

func NewHuddleHandler(service ports.HuddleService) *HuddleHandler {
	return &HuddleHandler{service: service}
}

// StartHuddle handles POST /huddle/start
func (h *HuddleHandler) StartHuddle(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		ChannelID string `json:"channel_id"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	huddle, err := h.service.StartHuddle(r.Context(), body.ChannelID)
	if err != nil {
		http.Error(w, "Failed to create huddle: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(huddle)
}

// LogActivity handles POST /huddle/log
func (h *HuddleHandler) LogActivity(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		HuddleID  string         `json:"huddle_id"`
		UserID    string         `json:"user_id"`
		EventType string         `json:"event_type"`
		MetaData  map[string]any `json:"meta_data"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.service.LogActivity(r.Context(), body.HuddleID, body.UserID, body.EventType, body.MetaData)
	if err != nil {
		http.Error(w, "Failed to log activity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
}

// GetHuddleSummary handles GET /huddle/{id}/summary
func (h *HuddleHandler) GetHuddleSummary(w http.ResponseWriter, r *http.Request) {
	// Split the path to get the ID (assumes /huddle/UUID/summary)
	parts := splitPath(r.URL.Path)

	// The path should look like ["huddle", "UUID", "summary"] (length 3)
	if len(parts) < 3 { 
		http.Error(w, "Missing huddle ID in path", http.StatusBadRequest)
		return
	}

	huddleID := parts[1] // The UUID should be at index 1

	summary, err := h.service.GetHuddleSummary(r.Context(), huddleID)
	if err != nil {
		http.Error(w, "Failed to retrieve summary: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// Helper function to safely split the path
func splitPath(path string) []string {
	var parts []string
	for _, part := range split(path, '/') {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// Helper function (since r.URL.Path doesn't guarantee leading/trailing slashes)
func split(s string, sep rune) []string {
	// Standard library split is fine for this context
	return strings.Split(s, string(sep))
}