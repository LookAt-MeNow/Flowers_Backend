package router

import (
	"fmt"
	"github.com/LookAt-MeNow/flowers/models"
	"github.com/LookAt-MeNow/flowers/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"errors"
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

		// 商品相关路由
		goods := api.Group("/goods")
		{
			// 搜索相关路由
			goods.GET("/qsearch", func(c *gin.Context) {
				searchHandler(c, db) // 将 db 传递给 searchHandler
			})
			//商品列表搜索
			goods.GET("/search", func(c *gin.Context) {
				search_goods_list(c, db) // 将 db 传递给 searchHandler
			})
			// 商品详情
			goods.GET("/detail", func(c *gin.Context) {
				goodsDetailHandler(c, db) // 将 db 传递给 goodsDetailHandler
			})

		}
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
	//指定goods_search表
	result := db.Table("goods_search"). // 指定表名为 goods
						Select("goods_id, goods_name").
						Where("goods_name LIKE ?", "%"+query+"%").
						Limit(6).
						Find(&goods)

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
		"message": goods,
	})
}

// 新增分页搜索处理函数
func search_goods_list(c *gin.Context, db *gorm.DB) {
	// 获取并验证参数
	query := c.Query("query")
	cid := c.Query("cid")
	pagenum, _ := strconv.Atoi(c.DefaultQuery("pagenum", "1"))
	pagesize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "10"))

	if pagenum < 1 {
		pagenum = 1
	}
	if pagesize < 1 || pagesize > 100 {
		pagesize = 10
	}
	offset := (pagenum - 1) * pagesize

	// 构建查询
	queryBuilder := db.Table("goods").Select("*")

	// 添加查询条件
	if query != "" {
		queryBuilder = queryBuilder.Where("goods_name LIKE ?", "%"+query+"%")
	}
	if cid != "" {
		queryBuilder = queryBuilder.Where("cat_id = ?", cid)
	}

	// 执行分页查询
	var goods []models.Goods
	result := queryBuilder.
		Offset(offset).
		Limit(pagesize).
		Find(&goods)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"msg":    "数据库查询失败",
				"status": http.StatusInternalServerError,
			},
		})
		return
	}

	// 获取总记录数
	var total int64
	countQuery := db.Table("goods") // 指定表名为 goods
	if query != "" {
		countQuery = countQuery.Where("goods_name LIKE ?", "%"+query+"%")
	}
	if cid != "" {
		countQuery = countQuery.Where("cat_id = ?", cid)
	}
	countQuery.Count(&total)

	// 构建响应
	response := gin.H{
		"total":   total,
		"pagenum": pagenum,
		"goods":   goods,
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Message: response,
		Meta: models.Meta{
			Msg:    "获取成功",
			Status: http.StatusOK,
		},
	})
}

// 商品详情
func goodsDetailHandler(c *gin.Context, db *gorm.DB) {
	goodsID := c.Query("goods_id")
	if goodsID == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "商品ID不能为空",
				Status: http.StatusBadRequest,
			},
		})
		return
	}

	// 查询 goods 表
	var goods models.Goods
	result := db.Where("goods_id = ?", goodsID).First(&goods)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ApiResponse{
				Meta: models.Meta{
					Msg:    "商品不存在",
					Status: http.StatusNotFound,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Meta: models.Meta{
					Msg:    "服务器错误",
					Status: http.StatusInternalServerError,
				},
			})
		}
		return
	}

	// 查询 goods_detail 表
	var goodsDetail models.Goods_detail
	result = db.Where("goods_id = ?", goodsID).First(&goodsDetail)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "获取商品详情失败",
				Status: http.StatusInternalServerError,
			},
		})
		return
	}

	// 查询商品图片和属性
	var pics []models.GoodsPicture
	db.Where("goods_id = ?", goodsID).Find(&pics)

	var attrs []models.GoodsAttr
	db.Where("goods_id = ?", goodsID).Find(&attrs)

	// 构建响应数据结构
	type ResponseDetail struct {
		models.Goods
		GoodsIntroduce string                `json:"goods_introduce"`
		GoodsState     int                   `json:"goods_state"`
		IsDel          string                `json:"is_del"`
		Pics           []models.GoodsPicture `json:"pics"`
		Attrs          []models.GoodsAttr    `json:"attrs"`
		AddTime        int64                 `json:"add_time"`
		UpdTime        int64                 `json:"upd_time"`
	}

	response := ResponseDetail{
		Goods:          goods,
		GoodsIntroduce: goodsDetail.GoodsIntroduce,
		GoodsState:     goodsDetail.GoodsState,
		IsDel:          goodsDetail.IsDel,
		Pics:           pics,
		Attrs:          attrs,
		AddTime:        goods.AddTime.Unix(),
		UpdTime:        goods.UpdTime.Unix(),
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Message: response,
		Meta: models.Meta{
			Msg:    "获取成功",
			Status: http.StatusOK,
		},
	})
}
