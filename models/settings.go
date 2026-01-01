package models

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Setting struct {
	ID    uint           `gorm:"primaryKey"`
	Key   string         `gorm:"uniqueIndex;size:255"`
	Value datatypes.JSON `gorm:"type:jsonb"` // use type:json for MySQL
}

// GetSetting fetches a setting by key
func GetSetting(db *gorm.DB, key string, out any) (datatypes.JSON, error) {
	var setting Setting
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return setting.Value, nil
}

// SetSetting creates or updates a JSON setting
func SetSetting(db *gorm.DB, key string, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	setting := Setting{
		Key:   key,
		Value: jsonValue,
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		UpdateAll: true,
	}).Create(&setting).Error
}

func GetAllSettingsAsMap(db *gorm.DB) (map[string]map[string]any, error) {
	var settings []Setting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]map[string]any, len(settings))
	for _, s := range settings {
		var value map[string]any
		if err := json.Unmarshal(s.Value, &value); err != nil {
			return nil, err
		}
		result[s.Key] = value
	}
	return result, nil
}
