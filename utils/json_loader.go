package utils

import (
    "encoding/json"
    "os"
    //"path/filepath"
	//"net/http"
)

// 从本地加载数据 interface{}是空接口，可以接受任何类型的数据
/* func LoadJSONData(path string, target interface{}) error {
    absPath, _ := filepath.Abs(path)
    file, err := os.ReadFile(absPath)
    if err != nil {
        return err
    }
    return json.Unmarshal(file, target)
} */

func LoadJSONData(filePath string, target interface{}) error {
    // 读取文件
    file, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    // 解析 JSON 数据
    return json.Unmarshal(file, target) //unmarshal是将json数据解析为结构体
}

/* func SendJSON(w http.ResponseWriter, data interface{}) {
	// 设置响应头为 JSON 格式
    w.Header().Set("Content-Type", "application/json")
	// 返回 JSON 数据
    json.NewEncoder(w).Encode(data)
} */