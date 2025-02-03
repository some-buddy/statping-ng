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

	// Récupérer la configuration des statuts intermédiaires depuis la base de données
	err := db.First(&outage).Error
	if err != nil {
		// Si aucune configuration n'est trouvée, retourner une erreur
		http.Error(w, "Failed to fetch outage configuration", http.StatusInternalServerError)
		return
	}

	// Retourner la configuration en format JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(outage); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Handler pour enregistrer ou mettre à jour la configuration des statuts intermédiaires
func OutageConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var outage outage.OutageConfig

	// Lire les données envoyées dans le corps de la requête (en JSON)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&outage); err != nil {
		// Si le décodage échoue, retourner une erreur
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Vérifier si une configuration existe déjà
	var existingOutage outage.OutageConfig
	err := db.First(&existingOutage).Error
	if err != nil {
		// Si aucune configuration n'existe (pour POST), on la crée
		err = db.Create(&outage).Error
		if err != nil {
			http.Error(w, "Failed to create outage configuration", http.StatusInternalServerError)
			return
		}
	} else {
		// Sinon (pour PUT), on met à jour la configuration existante
		err = db.Model(&existingOutage).Updates(outage).Error
		if err != nil {
			http.Error(w, "Failed to update outage configuration", http.StatusInternalServerError)
			return
		}
	}

	// Répondre avec un message de succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Configuration saved successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
