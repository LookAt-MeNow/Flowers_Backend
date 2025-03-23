package models

//分类数据结构
type Category struct {
    Name         string `json:"name"`
    ImageSrc     string `json:"image_src"`
    OpenType     string `json:"open_type,omitempty"`
    NavigatorURL string `json:"navigator_url,omitempty"`
}

////分类页面结构
type CategoryTree struct {
    CatID      int           `json:"cat_id"`
    CatName    string        `json:"cat_name"`
    CatPid     int           `json:"cat_pid"`
    CatLevel   int           `json:"cat_level"`
    CatDeleted bool          `json:"cat_deleted"`
    CatIcon    string        `json:"cat_icon"`
    Children   []CategoryTree `json:"children,omitempty"`
}