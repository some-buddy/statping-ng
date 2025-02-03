package outage

import (
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/errors"  // For error handling
	"fmt" // For string formatting
)

var (
	allOutages map[int64]*OutageConfig
)

// SetDB initializes the db variable with a database.Database instance
func SetDB(database database.Database) {
	db = database // Set the db to the provided database.Database instance
}

// DB returns the current database.Database instance
func DB() database.Database {
	return db
}

// All retrieves all OutageConfig records from the database
func All() map[int64]*OutageConfig {
	return allOutages
}

// Find retrieves a specific OutageConfig by ID
func Find(id int64) (*OutageConfig, error) {
	outage := allOutages[id]
	if outage == nil {
		return nil, errors.Missing(&OutageConfig{}, id)
	}
	db.First(&outage, id)
	return outage, nil
}

// Create adds a new OutageConfig to the database
func (i *OutageConfig) Create() error {
	err := db.Create(i)
	if err.Error() != nil {
		log.Errorln(fmt.Sprintf("Failed to create outage #%v: %v", i.Id, err))
		return err.Error()
	}
	return nil
}

// Update updates an existing OutageConfig in the database
func (i *OutageConfig) Update() error {
	q := db.Update(i)
	allOutages[i.Id] = i
	return q.Error()
}

// Delete removes the OutageConfig from the database
func (i *OutageConfig) Delete() error {
	delete(allOutages, i.Id)
	q := db.Model(&OutageConfig{}).Delete(i)
	return q.Error()
}