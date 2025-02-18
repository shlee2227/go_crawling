package repository

import (
	entities "github.com/shlee2227/go_crawling/internal/entities/search"
	"gorm.io/gorm"
)


type Repository interface { 
 	Create([]entities.Item) error 
	GetAll() ([]entities.Item, error)
}

type repository struct {
	db *gorm.DB
}

// Repository 객체 생성 
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db, 
	}
}

// interface 기반 method 구현
func (r *repository) Create(Items []entities.Item) error {
	return r.db.Create(&Items).Error
} 

func (r *repository) GetAll() ([]entities.Item, error) {
    var items []entities.Item
    err := r.db.Find(&items).Error
    return items, err
}