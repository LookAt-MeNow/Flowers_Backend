package models

import (
	"time"
)
// Goods 商品数据结构
type Goods struct {
	GoodsID         uint      `gorm:"primaryKey;column:goods_id" json:"goods_id"`
	CatID           uint      `gorm:"column:cat_id" json:"cat_id"`
	GoodsName       string    `gorm:"column:goods_name" json:"goods_name"`
	GoodsPrice      float64   `gorm:"type:decimal(10,2);column:goods_price" json:"goods_price"`
	GoodsNumber     uint      `gorm:"column:goods_number" json:"goods_number"`
	GoodsWeight     uint      `gorm:"column:goods_weight" json:"goods_weight"`
	GoodsBigLogo    string    `gorm:"column:goods_big_logo" json:"goods_big_logo"`
	GoodsSmallLogo  string    `gorm:"column:goods_small_logo" json:"goods_small_logo"`
	AddTime         time.Time `gorm:"autoCreateTime;column:add_time" json:"add_time"`
	UpdTime         time.Time `gorm:"autoUpdateTime;column:upd_time" json:"upd_time"`
	IsPromote       bool      `gorm:"type:tinyint(1);column:is_promote" json:"is_promote"`
	HotNumber       uint      `gorm:"column:hot_number" json:"hot_number"`
}

//------------------------------------------------------------------------
// Goods 商品数据结构
type Goods_search struct {
	ID       uint   `gorm:"primaryKey;column:goods_id"`
	Name     string `gorm:"column:goods_name"`
	// 添加其他字段...
}

//----------------------------------------------------------------------
// Goods_detail数据结构
// 商品详情模型（继承Goods并扩展）
type Goods_detail struct {
	Goods                  // 内嵌基础商品结构体
	GoodsIntroduce  string `gorm:"type:text;column:goods_introduce" json:"goods_introduce"`
	GoodsState      int    `gorm:"column:goods_state" json:"goods_state"`
	IsDel           string `gorm:"column:is_del" json:"is_del"`
	
	// 关联模型
	Pics    []GoodsPicture   `gorm:"foreignKey:GoodsID" json:"pics"`
	Attrs   []GoodsAttr      `gorm:"foreignKey:GoodsID" json:"attrs"`
}

// 自定义表名
func (Goods_detail) TableName() string {
	return "goods_detail"
}

// 商品图片模型
type GoodsPicture struct {
	PicsID    uint   `gorm:"primaryKey;column:pics_id" json:"pics_id"`
	GoodsID   uint   `gorm:"index;column:goods_id" json:"goods_id"`
	PicsBig   string `gorm:"column:pics_big" json:"pics_big"`
	PicsMid   string `gorm:"column:pics_mid" json:"pics_mid"`
	PicsSma   string `gorm:"column:pics_sma" json:"pics_sma"`
}

// 商品属性模型
type GoodsAttr struct {
	AttrID     uint    `gorm:"primaryKey;column:attr_id" json:"attr_id"`
	GoodsID    uint    `gorm:"index;column:goods_id" json:"goods_id"`
	AttrValue  string  `gorm:"column:attr_value" json:"attr_value"`
	AddPrice   float64 `gorm:"type:decimal(10,2);column:add_price" json:"add_price"`
	AttrName   string  `gorm:"column:attr_name" json:"attr_name"`
	AttrSel    string  `gorm:"column:attr_sel" json:"attr_sel"`
	AttrWrite  string  `gorm:"column:attr_write" json:"attr_write"`
	AttrVals   string  `gorm:"column:attr_vals" json:"attr_vals"`
}