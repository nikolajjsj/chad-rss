package database

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	Close() error

	// User CRUD
	GetUserByUsername(username string) (User, error)
	GetUserByID(id uint) (User, error)
	CreateUser(username string, password string) (User, error)
	// Feeds & Articles CRUD
	GetAllFeeds() ([]Feed, error)
}

type service struct {
	DB *gorm.DB
}

var (
	database   = os.Getenv("DATABASE")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Feed{}, &Article{})

	dbInstance = &service{
		DB: db,
	}
	return dbInstance
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB.Close()
}

/// Query methods

// Get user by username
func (s *service) GetUserByUsername(username string) (User, error) {
	var user User
	result := s.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

// Get user by ID
func (s *service) GetUserByID(id uint) (User, error) {
	var user User
	result := s.DB.First(&user, id)
	return user, result.Error
}

// Create user
func (s *service) CreateUser(username string, password string) (User, error) {
	user := User{
		Username: username,
		Password: password,
	}
	result := s.DB.Create(&user)
	return user, result.Error
}

// Get all feeds
func (s *service) GetAllFeeds() ([]Feed, error) {
	var feeds []Feed
	result := s.DB.Find(&feeds)
	return feeds, result.Error
}
