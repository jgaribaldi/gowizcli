package db

import (
	"errors"
	"fmt"
	"gowizcli/wiz"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage interface {
	Upsert(bulb wiz.Light) (*wiz.Light, error)
	FindAll() ([]wiz.Light, error)
	EraseAll()
	FindById(id string) (*wiz.Light, error)
}

type SQLiteDB struct {
	db *gorm.DB
}

func NewSQLiteDB(filename string) (*SQLiteDB, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&storedWizLight{})

	return &SQLiteDB{db: db}, nil
}

func (s SQLiteDB) Upsert(bulb wiz.Light) (*wiz.Light, error) {
	storedWizLight := storedWizLight{
		ID:         bulb.Id,
		MacAddress: bulb.MacAddress,
		IpAddress:  bulb.IpAddress,
	}

	result := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mac_address"}},
		DoUpdates: clause.AssignmentColumns([]string{"ip_address"}),
	}).Create(&storedWizLight)
	if result.Error != nil {
		return nil, result.Error
	}

	return &wiz.Light{
		Id:         storedWizLight.ID,
		MacAddress: storedWizLight.MacAddress,
		IpAddress:  storedWizLight.IpAddress,
	}, nil
}

func (s SQLiteDB) FindAll() ([]wiz.Light, error) {
	var storedWizLights []storedWizLight
	storedWizLights = make([]storedWizLight, 0)
	queryResult := s.db.Find(&storedWizLights)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	var result []wiz.Light
	result = make([]wiz.Light, 0)
	for _, l := range storedWizLights {
		result = append(result, wiz.Light{
			Id:         l.ID,
			IpAddress:  l.IpAddress,
			MacAddress: l.MacAddress,
		})
	}
	return result, nil
}

func (s SQLiteDB) EraseAll() {
	tableName := "stored_lights"
	s.db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
}

func (s SQLiteDB) FindById(id string) (*wiz.Light, error) {
	storedWizLight := storedWizLight{ID: id}

	queryResult := s.db.First(&storedWizLight)
	if queryResult.Error != nil && errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("id %s not found", id)
	}
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return &wiz.Light{
		Id:         storedWizLight.ID,
		IpAddress:  storedWizLight.IpAddress,
		MacAddress: storedWizLight.MacAddress,
	}, nil
}

type storedWizLight struct {
	gorm.Model
	ID         string
	MacAddress string `gorm:"uniqueIndex"`
	IpAddress  string `gorm:"uniqueIndex"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (storedWizLight) TableName() string {
	return "stored_lights"
}
