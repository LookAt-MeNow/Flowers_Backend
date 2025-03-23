/* package router

import (
	"encoding/json"
    "net/http"
    "github.com/LookAt-MeNow/flowers/models"
    "github.com/LookAt-MeNow/flowers/utils"
	"fmt"
)

func RegisterRoutes() {
	// 注册路由
	http.HandleFunc("/api/public/v1/home/swiperdata", swiperHandler) // 轮播图路由
	http.HandleFunc("/api/public/v1/home/catitems", catItemsHandler) // 分类路由
	http.HandleFunc("/api/public/v1/home/floordata", floorHandler)   // 楼层路由
	http.HandleFunc("/api/public/v1/categories", categoriesHandler) // 分类路由
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
    response := models.ApiResponse{
        Message: data,
        Meta: models.Meta{
            Msg:    "获取成功",
            Status: http.StatusOK,
        },
    }
	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")
	// 返回 JSON 数据
	json.NewEncoder(w).Encode(response)
    //utils.SendJSON(w, response)
}


// 轮播图数据处理
func swiperHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("接受到轮播图请求")
    var banners []models.Banner
    if err := utils.LoadJSONData("data/swiperdata.json", &banners); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    sendJSONResponse(w, banners)
}

// 分类页面数据处理
func catItemsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("接受到分类主页请求")
	var catItems []models.Category
	if err := utils.LoadJSONData("data/catitems.json", &catItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, catItems)
}

// 楼层数据处理
func floorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("接受到楼层请求")
	var floors []models.Floor
	if err := utils.LoadJSONData("data/floordata.json", &floors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, floors)
}

// 分类数据处理
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("接受到分类页面请求")
	var categories []models.CategoryTree
	if err := utils.LoadJSONData("data/categories.json", &categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, categories)
}

 */

 package router

 import (
	 "net/http"
    "github.com/LookAt-MeNow/flowers/models"
    "github.com/LookAt-MeNow/flowers/utils"
	 "github.com/gin-gonic/gin"
	 "gorm.io/gorm"
	 "fmt"
 )
 
 // SetupRouter 初始化 Gin 路由并返回引擎实例
 func SetupRouter(db *gorm.DB) *gin.Engine {
	 r := gin.Default()
 
	 // 配置公共中间件
	 r.Use(CORSMiddleware())
	 r.Use(ResponseWrapper())
 
	 // 注册路由组
	 api := r.Group("/api/public/v1")
	 {
		 // 首页相关路由
		 home := api.Group("/home")
		 {
			 home.GET("/swiperdata", swiperHandler)
			 home.GET("/catitems", catItemsHandler)
			 home.GET("/floordata", floorHandler)
		 }
 
		 // 分类相关路由
		 api.GET("/categories", categoriesHandler)
 
		 // 搜索相关路由
		 api.GET("/goods/qsearch", func(c *gin.Context) {
			 searchHandler(c, db) // 将 db 传递给 searchHandler
		 })
	 }
	 return r
 }
 
 // CORSMiddleware 跨域中间件
 func CORSMiddleware() gin.HandlerFunc {
	 return func(c *gin.Context) {
		 c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		 c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		 c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		 if c.Request.Method == "OPTIONS" {
			 c.AbortWithStatus(204)
			 return
		 }
		 c.Next()
	 }
 }
 
 // ResponseWrapper 统一响应格式中间件
 func ResponseWrapper() gin.HandlerFunc {
	 return func(c *gin.Context) {
		 c.Next()
 
		 // 跳过静态文件请求
		 if c.Writer.Status() == http.StatusNotFound && c.Request.URL.Path == "/" {
			 return
		 }
 
		 // 统一处理响应格式
		 if response, exists := c.Get("response"); exists {
			 c.JSON(c.Writer.Status(), models.ApiResponse{
				 Message: response,
				 Meta: models.Meta{
					 Msg:    "获取成功",
					 Status: c.Writer.Status(),
				 },
			 })
		 }
	 }
 }
 
 // 轮播图数据处理
 func swiperHandler(c *gin.Context) {
	 var banners []models.Banner
	 if err := utils.LoadJSONData("data/swiperdata.json", &banners); err != nil {
		 c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			 "error": "轮播图数据加载失败",
		 })
		 return
	 }
	 c.Set("response", banners)
 }
 
 // 分类页面数据处理
 func catItemsHandler(c *gin.Context) {
	 var catItems []models.Category
	 if err := utils.LoadJSONData("data/catitems.json", &catItems); err != nil {
		 c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			 "error": "分类数据加载失败",
		 })
		 return
	 }
	 c.Set("response", catItems)
 }
 
 // 楼层数据处理
 func floorHandler(c *gin.Context) {
	 var floors []models.Floor
	 if err := utils.LoadJSONData("data/floordata.json", &floors); err != nil {
		 c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			 "error": "楼层数据加载失败",
		 })
		 return
	 }
	 c.Set("response", floors)
 }
 
 // 分类数据处理
 func categoriesHandler(c *gin.Context) {
	 var categories []models.CategoryTree
	 if err := utils.LoadJSONData("data/categories.json", &categories); err != nil {
		 c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			 "error": "分类树数据加载失败",
		 })
		 return
	 }
	 c.Set("response", categories)
 }
 
 // 搜索数据处理
 func searchHandler(c *gin.Context, db *gorm.DB) {
	 // 获取查询参数
	 query := c.Query("query")
	 fmt.Println("Received query:", query) // 打印查询参数
	 if query == "" {
		 c.JSON(http.StatusBadRequest, gin.H{
			 "error": "搜索关键字不能为空",
		 })
		 return
	 }
 
	 // 数据库查询
	 var goods []models.Goods_search // 用于保存查询结果
	 result := db.Select("goods_id, goods_name").
		 Where("goods_name LIKE ?", "%"+query+"%").
		 Limit(6).    // 限制查询结果数量
		 Find(&goods) // 查询结果保存到 goods 中
 
	 // 处理查询结果
	 if result.Error != nil {
		 c.JSON(http.StatusInternalServerError, gin.H{
			 "error": "数据库查询失败",
		 })
		 return
	 }
 
	 // 返回标准化响应
	 c.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"msg":    "获取成功",
			"status": http.StatusOK,
		},
		"message":  goods,
	 })
 }