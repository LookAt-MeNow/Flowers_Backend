package models

// Goods 商品数据结构
type Goods_search struct {
	ID       uint   `gorm:"primaryKey;column:goods_id"`
	Name     string `gorm:"column:goods_name"`
	// 添加其他字段...
}

func (Goods_search) TableName() string {
	return "goods"
}