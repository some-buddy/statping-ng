package intermediate_status

import (
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/errors"  // For error handling
	"fmt" // For string formatting
)

var (
	allIntermediateStatuses map[int64]*IntermediateStatusConfig
)

// SetDB initializes the db variable with a database.Database instance
func SetDB(database database.Database) {
	db = database // Set the db to the provided database.Database instance
}

// DB returns the current database.Database instance
func DB() database.Database {
	return db
}

// All retrieves all IntermediateStatusConfig records from the database
func All() map[int64]*IntermediateStatusConfig {
	return allIntermediateStatuses
}

// Find retrieves a specific IntermediateStatusConfig by ID
func Find(id int64) (*IntermediateStatusConfig, error) {
	status := allIntermediateStatuses[id]
	if status == nil {
		return nil, errors.Missing(&IntermediateStatusConfig{}, id)
	}
	db.First(&status, id)
	return status, nil
}

// Create adds a new IntermediateStatusConfig to the database
func (i *IntermediateStatusConfig) Create() error {
	err := db.Create(i)
	if err.Error() != nil {
		log.Errorln(fmt.Sprintf("Failed to create intermediate status #%v: %v", i.Id, err))
		return err.Error()
	}
	return nil
}

// Update updates an existing IntermediateStatusConfig in the database
func (i *IntermediateStatusConfig) Update() error {
	q := db.Update(i)
	allIntermediateStatuses[i.Id] = i
	return q.Error()
}

// Delete removes the IntermediateStatusConfig from the database
func (i *IntermediateStatusConfig) Delete() error {
	delete(allIntermediateStatuses, i.Id)
	q := db.Model(&IntermediateStatusConfig{}).Delete(i)
	return q.Error()
}