package models

// 轮播图数据结构
type Banner struct { 
    ImageSrc     string `json:"image_src"`
    OpenType     string `json:"open_type"`
    GoodsID      int    `json:"goods_id"`
    NavigatorURL string `json:"navigator_url"`
}