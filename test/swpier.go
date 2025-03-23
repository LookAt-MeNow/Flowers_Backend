package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// 轮播图数据结构
type Banner struct {
	ImageSrc     string `json:"image_src"`
	OpenType     string `json:"open_type"`
	GoodsID      int    `json:"goods_id"`
	NavigatorURL string `json:"navigator_url"`
}

// 分类数据结构
type Category struct {
	Name         string `json:"name"`
	ImageSrc     string `json:"image_src"`
	OpenType     string `json:"open_type,omitempty"`      // omitempty 表示字段为空时不显示
	NavigatorURL string `json:"navigator_url,omitempty"` // omitempty 表示字段为空时不显示
}

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

//分类页面结构
type Category_c struct {
	CatID      int        `json:"cat_id"`
	CatName    string     `json:"cat_name"`
	CatPid     int        `json:"cat_pid"`
	CatLevel   int        `json:"cat_level"`
	CatDeleted bool       `json:"cat_deleted"`
	CatIcon    string     `json:"cat_icon"`
	Children   []Category_c `json:"children,omitempty"`
}

// 响应数据结构
type ApiResponse struct {
	Message interface{} `json:"message"` // 使用 interface{} 支持多种数据类型
	Meta    Meta        `json:"meta"`
}

// Meta 数据结构
type Meta struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

// 定义轮播图数据
var banners = []Banner{
	{
		ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/banner1.png",
		OpenType:     "navigate",
		GoodsID:      129,
		NavigatorURL: "/pages/goods_detail/main?goods_id=129",
	},
	{
		ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/banner2.png",
		OpenType:     "navigate",
		GoodsID:      395,
		NavigatorURL: "/pages/goods_detail/main?goods_id=395",
	},
	{
		ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/banner3.png",
		OpenType:     "navigate",
		GoodsID:      38,
		NavigatorURL: "/pages/goods_detail/main?goods_id=38",
	},
}

// 定义分类数据
var categories = []Category{
	{
		Name:         "分类",
		ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/icon_index_nav_4@2x.png",
		OpenType:     "switchTab",
		NavigatorURL: "/pages/category/main",
	},
	{
		Name:     "秒杀拍",
		ImageSrc: "https://api-hmugo-web.itheima.net/pyg/icon_index_nav_3@2x.png",
	},
	{
		Name:     "超市购",
		ImageSrc: "https://api-hmugo-web.itheima.net/pyg/icon_index_nav_2@2x.png",
	},
	{
		Name:     "母婴品",
		ImageSrc: "https://api-hmugo-web.itheima.net/pyg/icon_index_nav_1@2x.png",
	},
}

// 定义楼层数据
var floors = []Floor{
	{
		FloorTitle: FloorTitle{
			Name:     "时尚女装",
			ImageSrc: "https://api-hmugo-web.itheima.net/pyg/pic_floor01_title.png",
		},
		ProductList: []FloorProduct{
			{
				Name:         "优质服饰",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor01_1@2x.png",
				ImageWidth:   "232",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=服饰",
			},
			{
				Name:         "春季热门",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor01_2@2x.png",
				ImageWidth:   "233",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=热",
			},
			{
				Name:         "爆款清仓",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor01_3@2x.png",
				ImageWidth:   "233",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=爆款",
			},
			{
				Name:         "倒春寒",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor01_4@2x.png",
				ImageWidth:   "233",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=春季",
			},
			{
				Name:         "怦然心动",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor01_5@2x.png",
				ImageWidth:   "233",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=心动",
			},
		},
	},
	{
		FloorTitle: FloorTitle{
			Name:     "户外活动",
			ImageSrc: "https://api-hmugo-web.itheima.net/pyg/pic_floor02_title.png",
		},
		ProductList: []FloorProduct{
			{
				Name:         "勇往直前",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor02_1@2x.png",
				ImageWidth:   "232",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=户外",
			},
			{
				Name:         "户外登山包",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor02_2@2x.png",
				ImageWidth:   "273",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=登山包",
			},
			{
				Name:         "超强手套",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor02_3@2x.png",
				ImageWidth:   "193",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=手套",
			},
			{
				Name:         "户外运动鞋",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor02_4@2x.png",
				ImageWidth:   "193",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=运动鞋",
			},
			{
				Name:         "冲锋衣系列",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor02_5@2x.png",
				ImageWidth:   "273",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=冲锋衣",
			},
		},
	},
	{
		FloorTitle: FloorTitle{
			Name:     "箱包配饰",
			ImageSrc: "https://api-hmugo-web.itheima.net/pyg/pic_floor03_title.png",
		},
		ProductList: []FloorProduct{
			{
				Name:         "清新气质",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor03_1@2x.png",
				ImageWidth:   "232",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=饰品",
			},
			{
				Name:         "复古胸针",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor03_2@2x.png",
				ImageWidth:   "263",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=胸针",
			},
			{
				Name:         "韩版手链",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor03_3@2x.png",
				ImageWidth:   "203",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=手链",
			},
			{
				Name:         "水晶项链",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor03_4@2x.png",
				ImageWidth:   "193",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=水晶项链",
			},
			{
				Name:         "情侣表",
				ImageSrc:     "https://api-hmugo-web.itheima.net/pyg/pic_floor03_5@2x.png",
				ImageWidth:   "273",
				OpenType:     "navigate",
				NavigatorURL: "/pages/goods_list?query=情侣表",
			},
		},
	},
}

// 处理轮播图 GET 请求
func swiperHandler(w http.ResponseWriter, r *http.Request) {
	// 接受到请求
	fmt.Println("接受到轮播图请求")

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 构造响应数据
	response := ApiResponse{
		Message: banners,
		Meta: Meta{
			Msg:    "获取成功",
			Status: 200,
		},
	}

	// 返回 JSON 数据
	json.NewEncoder(w).Encode(response)
}

// 处理分类 GET 请求
func catItemsHandler(w http.ResponseWriter, r *http.Request) {
	// 接受到请求
	fmt.Println("接受到分类请求")

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 构造响应数据
	response := ApiResponse{
		Message: categories,
		Meta: Meta{
			Msg:    "获取成功",
			Status: 200,
		},
	}

	// 返回 JSON 数据
	json.NewEncoder(w).Encode(response)
}

// 处理楼层 GET 请求
func floorHandler(w http.ResponseWriter, r *http.Request) {
	// 接受到请求
	fmt.Println("接受到楼层请求")

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 构造响应数据
	response := ApiResponse{
		Message: floors,
		Meta: Meta{
			Msg:    "获取成功",
			Status: 200,
		},
	}

	// 返回 JSON 数据
	json.NewEncoder(w).Encode(response)
}


// 加载 JSON 数据
func loadCategories(filename string) ([]Category_c, error) {
	// 读取文件
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件: %v", err)
	}
    //输出日志

	// 解析 JSON 数据
	var categories []Category_c
	if err := json.Unmarshal(file, &categories); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	return categories, nil
}

//处理分类 GET 请求，从categories.json文件中读取数据
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 接受到请求
	fmt.Println("接受到分类2请求")
	// 加载 JSON 数据
	categories, err := loadCategories("categories.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 构造响应数据
	response := ApiResponse{
		Message: categories,
		Meta: Meta{
			Msg:    "获取成功",
			Status: 200,
		},
	}

	// 返回 JSON 数据
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 注册路由
	http.HandleFunc("/api/public/v1/home/swiperdata", swiperHandler) // 轮播图路由
	http.HandleFunc("/api/public/v1/home/catitems", catItemsHandler) // 分类路由
	http.HandleFunc("/api/public/v1/home/floordata", floorHandler)   // 楼层路由
    http.HandleFunc("/api/public/v1/categories", categoriesHandler) // 分类路由

	// 启动服务
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}