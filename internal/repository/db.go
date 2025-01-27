package repository

import (
	"fmt"

	"github.com/priince9381/irm_backend/internal/config"
	"github.com/priince9381/irm_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.SocialAccount{},
		&models.Post{},
		&models.Media{},
		&models.Analytics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &Database{DB: db}, nil
}

// User repository methods
func (d *Database) CreateUser(user *models.User) error {
	return d.DB.Create(user).Error
}

func (d *Database) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := d.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Social Account repository methods
func (d *Database) CreateSocialAccount(account *models.SocialAccount) error {
	return d.DB.Create(account).Error
}

func (d *Database) GetUserSocialAccounts(userID string) ([]models.SocialAccount, error) {
	var accounts []models.SocialAccount
	err := d.DB.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

// Post repository methods
func (d *Database) CreatePost(post *models.Post) error {
	return d.DB.Create(post).Error
}

func (d *Database) GetUserPosts(userID string) ([]models.Post, error) {
	var posts []models.Post
	err := d.DB.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func (d *Database) UpdatePost(post *models.Post) error {
	return d.DB.Save(post).Error
}

// Media repository methods
func (d *Database) CreateMedia(media *models.Media) error {
	return d.DB.Create(media).Error
}

func (d *Database) GetUserMedia(userID string) ([]models.Media, error) {
	var media []models.Media
	err := d.DB.Where("user_id = ?", userID).Find(&media).Error
	return media, err
}
