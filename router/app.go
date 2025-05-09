package router

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/LookAt-MeNow/flowers/models"
	"github.com/LookAt-MeNow/flowers/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		// 用户认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/merchants/register", func(c *gin.Context) {
				merchantRegisterHandler(c, db)
			})
			auth.POST("/merchants/login", func(c *gin.Context) {
				merchantLoginHandler(c, db)
			})
			auth.POST("/admin/login", func(c *gin.Context) {
				adminLoginHandler(c, db)
			})
		}

		// 在 SetupRouter 鲜花上架修改
		merchant := api.Group("/merchants")
		{
			// 鲜花管理
			merchant.GET("/flowers", merchantListFlowersHandler(db))       // 获取鲜花列表
			merchant.POST("/flowers", merchantAddFlowerHandler(db))        // 添加鲜花
			merchant.GET("/flowers/:id", merchantGetFlowerHandler(db))     // 获取单个鲜花
			merchant.PUT("/flowers/:id", merchantUpdateFlowerHandler(db))  // 更新鲜花
			merchant.PUT("/flowers/:id/status", merchantUpdateFlowerStatusHandler(db)) // 更新状态
		}
	}
	return r
}

	// 在main.go或路由设置中添加CORS中间件
	func CORSMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

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


// --------------------------------------商家管理员端
// 商家注册处理
// 商家注册处理（简化版）
func merchantRegisterHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		ShopName string `json:"shop_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "参数错误",
				Status: http.StatusBadRequest,
			},
		})
		return
	}

	// 检查用户名是否已存在
	var count int64
	db.Model(&models.Merchant{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "用户名已存在",
				Status: http.StatusBadRequest,
			},
		})
		return
	}

	// 创建商家（密码明文存储，实际生产环境不推荐）
	merchant := models.Merchant{
		Username: req.Username,
		Password: req.Password, // 明文存储密码
		ShopName: req.ShopName,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1, // 默认激活状态
	}

	if err := db.Create(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "注册失败",
				Status: http.StatusInternalServerError,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Message: gin.H{
			"username": merchant.Username,
			"shopName": merchant.ShopName,
		},
		Meta: models.Meta{
			Msg:    "注册成功",
			Status: http.StatusCreated,
		},
	})
}

// 商家登录处理（简化版）
func merchantLoginHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "参数错误",
				Status: http.StatusBadRequest,
			},
		})
		return
	}

	// 查询商家
	var merchant models.Merchant
	if err := db.Where("username = ? AND password = ?", req.Username, req.Password).First(&merchant).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "用户名或密码错误",
				Status: http.StatusUnauthorized,
			},
		})
		return
	}

	// 检查商家状态
	if merchant.Status != 1 {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "商家账号已被禁用",
				Status: http.StatusForbidden,
			},
		})
		return
	}

	// 登录成功，返回商家信息
	c.JSON(http.StatusOK, models.ApiResponse{
		Message: gin.H{
			"merchant": gin.H{
				"id":       merchant.ID,
				"username": merchant.Username,
				"shopName": merchant.ShopName,
			},
		},
		Meta: models.Meta{
			Msg:    "登录成功",
			Status: http.StatusOK,
		},
	})
}

// 管理员登录处理（简化版）
func adminLoginHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "参数错误",
				Status: http.StatusBadRequest,
			},
		})
		return
	}

	// 查询管理员
	var admin models.Admin
	if err := db.Where("username = ? AND password = ?", req.Username, req.Password).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Meta: models.Meta{
				Msg:    "用户名或密码错误",
				Status: http.StatusUnauthorized,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Message: gin.H{
			"admin": gin.H{
				"id":       admin.ID,
				"username": admin.Username,
				"role":     admin.Role,
			},
		},
		Meta: models.Meta{
			Msg:    "登录成功",
			Status: http.StatusOK,
		},
	})
}

// 商家获取鲜花列表
func merchantListFlowersHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        merchantID := uint(1) // 临时使用固定值
        
        // 获取分页参数
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
        status := c.Query("status")
        search := c.Query("search")
        
        offset := (page - 1) * pageSize
        
        // 构建查询
        query := db.Model(&models.Flower{}).Where("merchant_id = ?", merchantID)
        
        if status != "" {
            query = query.Where("status = ?", status)
        }
        if search != "" {
            query = query.Where("name LIKE ?", "%"+search+"%")
        }
        
        // 获取总数
        var total int64
        if err := query.Count(&total).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "获取鲜花总数失败",
                "code": http.StatusInternalServerError,
            })
            return
        }
        
        // 获取分页数据
        var flowers []models.Flower
        if err := query.Offset(offset).Limit(pageSize).Preload("Images").Find(&flowers).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "获取鲜花列表失败",
                "code": http.StatusInternalServerError,
            })
            return
        }
        
        // 构建标准JSON响应
        response := gin.H{
            "data": gin.H{
                "list":  flowers,
                "total": total,
                "page":  page,
                "page_size": pageSize,
            },
            "code": http.StatusOK,
            "message": "获取成功",
        }
        
        c.JSON(http.StatusOK, response)
		//调试打印给前端的js数据
		fmt.Println("返回给前端的鲜花列表数据:", response)
    }
}

// 商家添加鲜花
func merchantAddFlowerHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        //merchantID := c.MustGet("merchantID").(uint)
		merchantID := uint(1) // 临时使用固定值 // 获取商家ID mustGet("merchantID") 是中间件设置的
            // 打印接收到的表单数据
        // 验证必填字段
        if c.PostForm("name") == ""  {
            c.JSON(http.StatusBadRequest, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "名称和价格不能为空",
                    Status: http.StatusBadRequest,
                },
            })
            return
        }
        
        // 处理图片上传
        form, err := c.MultipartForm()
        if err != nil {
            c.JSON(http.StatusBadRequest, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "请上传图片",
                    Status: http.StatusBadRequest,
                },
            })
            return
        }
        
        files := form.File["images"]
        if len(files) == 0 {
            c.JSON(http.StatusBadRequest, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "至少上传一张图片",
                    Status: http.StatusBadRequest,
                },
            })
            return
        }
    	// 解析表单数据
        name := c.PostForm("name")
        price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
        stock, _ := strconv.Atoi(c.PostForm("stock"))
        categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
        description := c.PostForm("description")
        status, _ := strconv.Atoi(c.PostForm("status"))

        // 保存鲜花基本信息
        flower := models.Flower{
            MerchantID:  merchantID,
            Name:        name,
            Price:       price,
            Stock:       stock,
            CategoryID:  uint(categoryID),
            Description: description,
            Status:      status,
        }
        
        if err := db.Create(&flower).Error; err != nil {
            c.JSON(http.StatusInternalServerError, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "创建鲜花失败",
                    Status: http.StatusInternalServerError,
                },
            })
            return
        }
        
        // 保存图片
        var imagePaths []string
        for _, file := range files {
            // 生成唯一文件名
            filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
            filepath := path.Join("uploads", filename)
            
            // 保存文件
            if err := c.SaveUploadedFile(file, filepath); err != nil {
                continue
            }
            
            // 保存到数据库
            image := models.FlowerImage{
                FlowerID: flower.ID,
                Path:     filepath,
            }
            if err := db.Create(&image).Error; err == nil {
                imagePaths = append(imagePaths, filepath)
            }
        }
        
        // 更新鲜花图片信息
        if len(imagePaths) > 0 {
			flower.Images = []models.FlowerImage{}
			for _, path := range imagePaths {
				flower.Images = append(flower.Images, models.FlowerImage{Path: path})
			}
            db.Save(&flower)
        }
        
        c.JSON(http.StatusCreated, models.ApiResponse{
            Message: flower,
            Meta: models.Meta{
                Msg:    "鲜花添加成功",
                Status: http.StatusCreated,
            },
        })
    }
}

// 商家获取单个鲜花详情
func merchantGetFlowerHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        //merchantID := c.MustGet("merchantID").(uint)
		merchantID := uint(1) // 临时使用固定值
        flowerID := c.Param("id")
        
        var flower models.Flower
        if err := db.Preload("Images").Where("id = ? AND merchant_id = ?", flowerID, merchantID).First(&flower).Error; err != nil {
            c.JSON(http.StatusNotFound, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "鲜花不存在或无权访问",
                    Status: http.StatusNotFound,
                },
            })
            return
        }
        
        c.JSON(http.StatusOK, models.ApiResponse{
            Message: flower,
            Meta: models.Meta{
                Msg:    "获取成功",
                Status: http.StatusOK,
            },
        })
    }
}

// 商家更新鲜花信息
func merchantUpdateFlowerHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        //merchantID := c.MustGet("merchantID").(uint)
		merchantID := uint(1) // 临时使用固定值
        flowerID := c.Param("id")
        
        var flower models.Flower
        if err := db.Preload("Images").Where("id = ? AND merchant_id = ?", flowerID, merchantID).First(&flower).Error; err != nil {
            c.JSON(http.StatusNotFound, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "鲜花不存在或无权访问",
                    Status: http.StatusNotFound,
                },
            })
            return
        }
        
        // 解析表单数据
        name := c.PostForm("name")
        price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
        stock, _ := strconv.Atoi(c.PostForm("stock"))
        categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
        description := c.PostForm("description")
        status, _ := strconv.Atoi(c.PostForm("status"))
        
        // 更新字段
        if name != "" {
            flower.Name = name
        }
        if price > 0 {
            flower.Price = price
        }
        if stock >= 0 {
            flower.Stock = stock
        }
        if categoryID > 0 {
            flower.CategoryID = uint(categoryID)
        }
        if description != "" {
            flower.Description = description
        }
        if status == 0 || status == 1 {
            flower.Status = status
        }
        
        // 处理图片上传
        form, err := c.MultipartForm()
        if err == nil {
            files := form.File["images"]
            var imagePaths []string
            
            // 先删除旧图片
            db.Where("flower_id = ?", flower.ID).Delete(&models.FlowerImage{})
            
            // 保存新图片
            for _, file := range files {
                filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
                filepath := path.Join("uploads", filename)
                
                if err := c.SaveUploadedFile(file, filepath); err != nil {
                    continue
                }
                
                image := models.FlowerImage{
                    FlowerID: flower.ID,
                    Path:     filepath,
                }
                if err := db.Create(&image).Error; err == nil {
                    imagePaths = append(imagePaths, filepath)
                }
            }
            
            if len(imagePaths) > 0 {
				flower.Images = []models.FlowerImage{}
				for _, path := range imagePaths {
					flower.Images = append(flower.Images, models.FlowerImage{Path: path})
				}
            }
        }
        
        if err := db.Save(&flower).Error; err != nil {
            c.JSON(http.StatusInternalServerError, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "更新鲜花失败",
                    Status: http.StatusInternalServerError,
                },
            })
            return
        }
        
        c.JSON(http.StatusOK, models.ApiResponse{
            Message: flower,
            Meta: models.Meta{
                Msg:    "鲜花更新成功",
                Status: http.StatusOK,
            },
        })
    }
}

// 商家更新鲜花状态
func merchantUpdateFlowerStatusHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        //merchantID := c.MustGet("merchantID").(uint)
		merchantID := uint(1) // 临时使用固定值
        flowerID := c.Param("id")
        
        var req struct {
            Status int `json:"status" binding:"required,oneof=0 1"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "参数错误",
                    Status: http.StatusBadRequest,
                },
            })
            return
        }
        
        var flower models.Flower
        if err := db.Where("id = ? AND merchant_id = ?", flowerID, merchantID).First(&flower).Error; err != nil {
            c.JSON(http.StatusNotFound, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "鲜花不存在或无权访问",
                    Status: http.StatusNotFound,
                },
            })
            return
        }
        
        flower.Status = req.Status
        if err := db.Save(&flower).Error; err != nil {
            c.JSON(http.StatusInternalServerError, models.ApiResponse{
                Meta: models.Meta{
                    Msg:    "更新状态失败",
                    Status: http.StatusInternalServerError,
                },
            })
            return
        }
        
        c.JSON(http.StatusOK, models.ApiResponse{
            Message: flower,
            Meta: models.Meta{
                Msg:    "状态更新成功",
                Status: http.StatusOK,
            },
        })
    }
}

