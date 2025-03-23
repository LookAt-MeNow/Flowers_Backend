package models

// 响应数据结构
type ApiResponse struct {
    Message interface{} `json:"message"`
    Meta    Meta        `json:"meta"`
}

type Meta struct {
    Msg    string `json:"msg"`
    Status int    `json:"status"`
}