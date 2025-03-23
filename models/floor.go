package models

// 楼层标题数据结构
type FloorTitle struct {
    Name     string `json:"name"`
    ImageSrc string `json:"image_src"`
}

// 楼层商品数据结构
type FloorProduct struct {
    Name         string `json:"name"`
    ImageSrc     string `json:"image_src"`
    ImageWidth   string `json:"image_width"`
    OpenType     string `json:"open_type"`
    NavigatorURL string `json:"navigator_url"`
}

// 楼层数据结构
type Floor struct {
    FloorTitle  FloorTitle     `json:"floor_title"`
    ProductList []FloorProduct `json:"product_list"`
}