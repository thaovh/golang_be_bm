package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterCommandRepo,
	NewGreeterQueryRepo,
	NewUserCommandRepo,
	NewUserQueryRepo,
	NewAuthCommandRepo,
	NewAuthQueryRepo,
	NewCountryCommandRepo,
	NewCountryQueryRepo,
	NewProvinceCommandRepo,
	NewProvinceQueryRepo,
	NewWardCommandRepo,
	NewWardQueryRepo,
)

// Data chứa cả read và write database
type Data struct {
	readDB  *gorm.DB // Database cho read operations
	writeDB *gorm.DB // Database cho write operations
}

// NewData tạo connections cho cả read và write database
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)

	// Kết nối Write Database (Master)
	writeDB, err := gorm.Open(postgres.Open(c.WriteDatabase.Source), &gorm.Config{})
	if err != nil {
		logHelper.Errorf("Failed to open write database: %v", err)
		return nil, nil, err
	}

	// Kết nối Read Database (Replica hoặc cùng database)
	readDB, err := gorm.Open(postgres.Open(c.ReadDatabase.Source), &gorm.Config{})
	if err != nil {
		logHelper.Errorf("Failed to open read database: %v", err)
		return nil, nil, err
	}

	// Test write connection
	writeSQLDB, err := writeDB.DB()
	if err != nil {
		logHelper.Errorf("Failed to get write database instance: %v", err)
		return nil, nil, err
	}
	if err := writeSQLDB.Ping(); err != nil {
		logHelper.Errorf("Failed to ping write database: %v", err)
		return nil, nil, err
	}

	// Test read connection
	readSQLDB, err := readDB.DB()
	if err != nil {
		logHelper.Errorf("Failed to get read database instance: %v", err)
		return nil, nil, err
	}
	if err := readSQLDB.Ping(); err != nil {
		logHelper.Errorf("Failed to ping read database: %v", err)
		return nil, nil, err
	}

	logHelper.Info("Write database connection established successfully")
	logHelper.Info("Read database connection established successfully")

	cleanup := func() {
		logHelper.Info("closing the data resources")

		// Close write database
		writeSQLDB, err := writeDB.DB()
		if err == nil {
			if err := writeSQLDB.Close(); err != nil {
				logHelper.Errorf("Failed to close write database: %v", err)
			} else {
				logHelper.Info("Write database connection closed")
			}
		}

		// Close read database
		readSQLDB, err := readDB.DB()
		if err == nil {
			if err := readSQLDB.Close(); err != nil {
				logHelper.Errorf("Failed to close read database: %v", err)
			} else {
				logHelper.Info("Read database connection closed")
			}
		}
	}

	return &Data{
		readDB:  readDB,
		writeDB: writeDB,
	}, cleanup, nil
}

// GetReadDB returns the read database instance
func (d *Data) GetReadDB() *gorm.DB {
	return d.readDB
}

// GetWriteDB returns the write database instance
func (d *Data) GetWriteDB() *gorm.DB {
	return d.writeDB
}
