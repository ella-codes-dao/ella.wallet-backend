package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	PublicKey string `gorm:"unique;not null" json:"publicKey"`
	Address   string `gorm:"not null" json:"address"`
}
