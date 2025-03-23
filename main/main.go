/* package main

import (
    "fmt"
    "net/http"
    "github.com/LookAt-MeNow/flowers/router"
)

func main() {
    router.RegisterRoutes()
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
} */

package main

import (
	"github.com/LookAt-MeNow/flowers/router"
    "github.com/LookAt-MeNow/flowers/sql"
)

func main() {
    // 加载配置
	cfg := sql.LoadConfig()
	// 初始化数据库
	db := sql.InitDB(cfg)
	// 初始化路由
	r := router.SetupRouter(db)
	// 启动服务
	r.Run(":8080") // 默认监听 0.0.0.0:8080
}