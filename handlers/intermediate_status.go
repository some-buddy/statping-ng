package handlers

import (
	"github.com/statping-ng/statping-ng/types/intermediate_status"
	"net/http"
	"github.com/statping-ng/statping-ng/database"
	"encoding/json"
)

var db database.Database

// GetIntermediateStatusConfig retrieves the intermediate status configuration.
func IntermediateStatusConfigViewHandler(w http.ResponseWriter, r *http.Request) {
	var intermediateStatus intermediate_status.IntermediateStatusConfig

	// Récupérer la configuration des statuts intermédiaires depuis la base de données
	err := db.First(&intermediateStatus).Error
	if err != nil {
		// Si aucune configuration n'est trouvée, retourner une erreur
		http.Error(w, "Failed to fetch intermediate status configuration", http.StatusInternalServerError)
		return
	}

	// Retourner la configuration en format JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(intermediateStatus); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Handler pour enregistrer ou mettre à jour la configuration des statuts intermédiaires
func IntermediateStatusConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var intermediateStatus intermediate_status.IntermediateStatusConfig

	// Lire les données envoyées dans le corps de la requête (en JSON)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&intermediateStatus); err != nil {
		// Si le décodage échoue, retourner une erreur
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Vérifier si une configuration existe déjà
	var existingStatus intermediate_status.IntermediateStatusConfig
	err := db.First(&existingStatus).Error
	if err != nil {
		// Si aucune configuration n'existe (pour POST), on la crée
		err = db.Create(&intermediateStatus).Error
		if err != nil {
			http.Error(w, "Failed to create intermediate status configuration", http.StatusInternalServerError)
			return
		}
	} else {
		// Sinon (pour PUT), on met à jour la configuration existante
		err = db.Model(&existingStatus).Updates(intermediateStatus).Error
		if err != nil {
			http.Error(w, "Failed to update intermediate status configuration", http.StatusInternalServerError)
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
