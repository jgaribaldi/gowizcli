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
	Upsert(bulb wiz.WizLight) (*wiz.WizLight, error)
	FindAll() ([]wiz.WizLight, error)
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

func (c Connection) Upsert(bulb wiz.WizLight) (*wiz.WizLight, error) {
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

	return &wiz.WizLight{
		Id:         storedWizLight.ID,
		MacAddress: storedWizLight.MacAddress,
		IpAddress:  storedWizLight.IpAddress,
	}, nil
}

func (c Connection) FindAll() ([]wiz.WizLight, error) {
	var storedWizLights []storedWizLight
	storedWizLights = make([]storedWizLight, 0)
	queryResult := c.db.Find(&storedWizLights)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	var result []wiz.WizLight
	result = make([]wiz.WizLight, 0)
	for _, l := range storedWizLights {
		result = append(result, wiz.WizLight{
			Id:         l.ID,
			IpAddress:  l.IpAddress,
			MacAddress: l.MacAddress,
		})
	}
	return result, nil
}

func (c Connection) EraseAll() {
	tableName := "stored_lights"
	c.db.Exec(fmt.Sprintf("DROP TABLE %s", tableName))
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
