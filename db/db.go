package db

import (
	"fmt"
	"gowizcli/wiz"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LightsDatabase interface {
	Upsert(bulb wiz.Light) (*wiz.Light, error)
	FindAll() ([]wiz.Light, error)
	EraseAll()
}

type Connection struct {
	db *gorm.DB
}

func NewConnection(filename string) (*Connection, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&storedWizLight{})

	return &Connection{db: db}, nil
}

func (c Connection) Upsert(bulb wiz.Light) (*wiz.Light, error) {
	storedWizLight := storedWizLight{
		ID:         bulb.Id,
		MacAddress: bulb.MacAddress,
		IpAddress:  bulb.IpAddress,
	}

	result := c.db.Clauses(clause.OnConflict{
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

func (c Connection) FindAll() ([]wiz.Light, error) {
	var storedWizLights []storedWizLight
	storedWizLights = make([]storedWizLight, 0)
	queryResult := c.db.Find(&storedWizLights)

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

func (c Connection) EraseAll() {
	tableName := "stored_lights"
	c.db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
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
