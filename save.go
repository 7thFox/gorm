package gorm

// Save creates or updates a DB entry for the given DataObject.
// Create/Update logic is based on the PK being > 0 for updates.
func (db *DatabaseConnection) Save(obj DataObject) error {
	pk := db.getPrimaryKeyValue(obj)
	if pk > 0 {
		return db.update(obj)
	}
	return db.insert(obj)
}
