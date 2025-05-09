package models
import "gorm.io/gorm"
// Flower 鲜花模型
type Flower struct {
    gorm.Model
    MerchantID  uint           `gorm:"index;not null"`
    Name        string         `gorm:"size:100;not null"`
    Price       float64        `gorm:"type:decimal(10,2);not null"`
    Stock       int            `gorm:"not null"`
    CategoryID  uint           `gorm:"index"`
    Description string         `gorm:"type:text"`
    Status      int            `gorm:"default:1"` // 1-上架, 0-下架
    Images      []FlowerImage  `gorm:"foreignKey:FlowerID" json:"images"`
}

// FlowerImage 鲜花图片模型
type FlowerImage struct {
    gorm.Model
    FlowerID uint   `gorm:"index;not null"`
    Path     string `gorm:"size:255;not null"`
}