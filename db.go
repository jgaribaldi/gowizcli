package main

import (
	"fmt"
	"gowizcli/wiz"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBConnection struct {
	db *gorm.DB
}

func NewDbConnection(filename string) (*DBConnection, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&storedWizLight{})

	return &DBConnection{db: db}, nil
}

func (d DBConnection) Upsert(bulb wiz.WizLight) (*wiz.WizLight, error) {
	storedWizLight := storedWizLight{
		ID:         bulb.Id,
		MacAddress: bulb.MacAddress,
		IpAddress:  bulb.IpAddress,
	}

	result := d.db.Clauses(clause.OnConflict{
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

func (d DBConnection) FindAll() ([]wiz.WizLight, error) {
	var storedWizLights []storedWizLight
	storedWizLights = make([]storedWizLight, 0)
	queryResult := d.db.Find(&storedWizLights)

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

func (d DBConnection) Reset() {
	tableName := "stored_lights"
	d.db.Exec(fmt.Sprintf("DROP TABLE %s", tableName))
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
