package outage

import (
	"fmt"
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/utils"
)

var log = utils.Log.WithField("type", "outage")
var db database.Database

// InitializeFromConfig initializes the outage configuration from the database.
func InitializeFromConfig() (*OutageConfig, error) {
	// Check if the db is initialized
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var outageConfig OutageConfig

	// Essayer de récupérer la configuration dans la table outage
	err := db.FirstOrCreate(&outageConfig, OutageConfig{Id: 1}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Outage: %w", err)
	}

	// Ensure correct handling of Id
	if outageConfig.Id == 0 {
		outageConfig.Id = 1 // Default or set based on logic
		err := db.Save(&outageConfig).Error
		if err != nil {
			return nil, fmt.Errorf("failed to save Outage: %w", err)
		}
	}

	// Returning the configuration
	return &outageConfig, nil
}
