package handlers

import (
	"github.com/statping-ng/statping-ng/types/outage"
	"net/http"
	"github.com/statping-ng/statping-ng/database"
	"encoding/json"
)

var db database.Database

// GetOutageConfig retrieves the outage configuration.
func OutageConfigViewHandler(w http.ResponseWriter, r *http.Request) {
	var outage outage.OutageConfig

	err := db.First(&outage).Error
	if err != nil {
		http.Error(w, "Failed to fetch outage configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(outage); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Handler to save or update the outage status configuration
func OutageConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var outage outage.OutageConfig

	// Read the data sent in the request body (in JSON)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&outage); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Check if a configuration already exists
	err := db.First(&outage).Error
	if err != nil {
		// If no configuration exists (for POST), create it
		err = db.Create(&outage).Error
		if err != nil {
			http.Error(w, "Failed to create outage configuration", http.StatusInternalServerError)
			return
		}
	} else {
			// Otherwise (for PUT), update the existing configuration
		err = db.Model(&outage).Updates(outage).Error
		if err != nil {
			http.Error(w, "Failed to update outage configuration", http.StatusInternalServerError)
			return
		}
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Configuration saved successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
