package intermediate_status

import (
	"fmt"
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/utils"
)

var log = utils.Log.WithField("type", "intermediate_status")
var db database.Database

// InitializeFromConfig initializes the intermediate status configuration from the database.
func InitializeFromConfig() (*IntermediateStatusConfig, error) {
	// Check if the db is initialized
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var intermediateStatusConfig IntermediateStatusConfig

	// Essayer de récupérer la configuration dans la table intermediate_status
	err := db.FirstOrCreate(&intermediateStatusConfig, IntermediateStatusConfig{Id: 1}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to initialize IntermediateStatus: %w", err)
	}

	// Ensure correct handling of Id
	if intermediateStatusConfig.Id == 0 {
		intermediateStatusConfig.Id = 1 // Default or set based on logic
		err := db.Save(&intermediateStatusConfig).Error
		if err != nil {
			return nil, fmt.Errorf("failed to save IntermediateStatus: %w", err)
		}
	}

	// Returning the configuration
	return &intermediateStatusConfig, nil
}
