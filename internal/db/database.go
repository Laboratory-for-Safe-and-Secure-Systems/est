package db

import (
    "fmt"
    "log"
    "reflect"

    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type DB struct {
    conn *gorm.DB
}

// NewDB creates a new DB instance and initializes the database connection (SQLite or PostgreSQL).
func NewDB(dbType string, dsn string) (*DB, error) {
    var db *gorm.DB
    var err error

    // Connect to the appropriate database based on the dbType
    switch dbType {
    case "postgres":
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            return nil, err
        }
    case "sqlite":
        db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
        if err != nil {
            return nil, err
        }
    default:
        return nil, fmt.Errorf("unsupported database type: %s", dbType)
    }

    database := &DB{conn: db}

    // Run migrations using reflection
    err = database.AutoMigrateWithReflection()
    if err != nil {
        return nil, err
    }

    log.Println("Database migration successful!")
    return database, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
    sqlDB, err := db.conn.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

// AutoMigrateWithReflection finds and registers all GORM models using reflection
func (db *DB) AutoMigrateWithReflection() error {
    // List of all model types
    modelTypes := []interface{}{
        Certificate{},
        CSR{},
        Revocation{},
        HTTPRequest{},
    }

    for _, model := range modelTypes {
        modelType := reflect.TypeOf(model)
        if modelType.Kind() == reflect.Ptr {
            modelType = modelType.Elem()
        }

        log.Printf("Migrating model: %s", modelType.Name())
        err := db.conn.AutoMigrate(model)
        if err != nil {
            return fmt.Errorf("failed to migrate model %s: %w", modelType.Name(), err)
        }
    }

    log.Println("Database schema migration completed for all models!")
    return nil
}

// Create inserts a new record into the database.
func (db *DB) Create(record interface{}) error {
    result := db.conn.Create(record)
    return result.Error
}

// Find retrieves a record from the database by primary key.
func (db *DB) Find(record interface{}, id uint) error {
    result := db.conn.First(record, id)
    return result.Error
}

// Update updates an existing record in the database.
func (db *DB) Update(record interface{}) error {
    result := db.conn.Save(record)
    return result.Error
}

// Delete removes a record from the database by primary key.
func (db *DB) Delete(record interface{}) error {
    result := db.conn.Delete(record)
    return result.Error
}
