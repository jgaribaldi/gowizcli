package db

import (
	"errors"
	"fmt"
	"gowizcli/wiz"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage interface {
	Upsert(bulb wiz.Light) (*wiz.Light, error)
	FindAll() ([]wiz.Light, error)
	EraseAll()
	FindById(id string) (*wiz.Light, error)
	AddTags(bulbs []wiz.Light, tags []string) ([]wiz.Light, error)
	RemoveTags(bulbs []wiz.Light, tags []string) ([]wiz.Light, error)
	FindByTags(tags []string) ([]wiz.Light, error)
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

	queryResult := s.db.Find(&storedWizLights)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	result := make([]wiz.Light, len(storedWizLights))
	for i, l := range storedWizLights {
		result[i] = wiz.Light{
			Id:         l.ID,
			IpAddress:  l.IpAddress,
			MacAddress: l.MacAddress,
			Tags:       l.Tags.Data(),
		}
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
		Tags:       storedWizLight.Tags.Data(),
	}, nil
}

func (s SQLiteDB) AddTags(bulbs []wiz.Light, tags []string) ([]wiz.Light, error) {
	result := make([]wiz.Light, len(bulbs))

	for i, b := range bulbs {
		newTags := add(b.Tags, tags)

		queryResult := s.db.
			Model(&storedWizLight{}).
			Where("id = ?", b.Id).
			Update("tags", datatypes.NewJSONType(newTags))
		if queryResult.Error != nil {
			return nil, queryResult.Error
		}

		result[i] = wiz.Light{
			Id:         b.Id,
			IpAddress:  b.IpAddress,
			MacAddress: b.MacAddress,
			IsOn:       b.IsOn,
			Tags:       newTags,
		}
	}

	return result, nil
}

func add(source []string, toAdd []string) []string {
	existingMap := make(map[string]struct{})
	for _, e := range source {
		existingMap[e] = struct{}{}
	}

	result := make([]string, 0)
	result = append(result, source...)
	for _, e := range toAdd {
		if _, exists := existingMap[e]; !exists {
			result = append(result, e)
		}
	}

	return result
}

func (s SQLiteDB) RemoveTags(bulbs []wiz.Light, tags []string) ([]wiz.Light, error) {
	result := make([]wiz.Light, len(bulbs))

	for i, b := range bulbs {
		newTags := filter(b.Tags, tags)

		queryResult := s.db.
			Model(&storedWizLight{}).
			Where("id = ?", b.Id).
			Update("tags", datatypes.NewJSONType(newTags))
		if queryResult.Error != nil {
			return nil, queryResult.Error
		}

		result[i] = wiz.Light{
			Id:         b.Id,
			IpAddress:  b.IpAddress,
			MacAddress: b.MacAddress,
			IsOn:       b.IsOn,
			Tags:       newTags,
		}
	}

	return result, nil
}

func filter(source []string, toRemove []string) []string {
	removeMap := make(map[string]struct{})
	for _, e := range toRemove {
		removeMap[e] = struct{}{}
	}

	result := make([]string, 0)
	for _, e := range source {
		if _, exists := removeMap[e]; !exists {
			result = append(result, e)
		}
	}
	return result
}

func (s SQLiteDB) FindByTags(tags []string) ([]wiz.Light, error) {
	var storedWizLights []storedWizLight

	tx := s.db.Model(&storedWizLight{})
	for _, t := range tags {
		tx = tx.Where("exists (select 1 from json_each(tags) where value = ?)", t)
	}

	if err := tx.Find(&storedWizLights).Error; err != nil {
		return nil, err
	}

	result := make([]wiz.Light, len(storedWizLights))
	for i, s := range storedWizLights {
		result[i] = wiz.Light{
			Id:         s.ID,
			IpAddress:  s.IpAddress,
			MacAddress: s.MacAddress,
			Tags:       s.Tags.Data(),
		}
	}

	return result, nil
}

type storedWizLight struct {
	gorm.Model
	ID         string
	MacAddress string `gorm:"uniqueIndex"`
	IpAddress  string `gorm:"uniqueIndex"`
	Tags       datatypes.JSONType[[]string]
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (storedWizLight) TableName() string {
	return "stored_lights"
}
