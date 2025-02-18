package entities

// gorm을 사용할 데이터는 gorm 반영해서 작성 필요
type Item struct {
    ID          uint   `gorm:"primaryKey"`
    Title       string `gorm:"not null"`
    Link        string `gorm:"not null"`
    Description string
    BloggerName string `gorm:"column:blogger_name"` // 컬럼명 다른 애들은 명시 
    BloggerLink string `gorm:"column:blogger_link"`
    PostDate    string `gorm:"column:post_date"`
}

func (Item) TableName() string {
	return "naver_items"
}

type NaverResponse struct {
    LastBuildDate string `json:"lastBuildDate"`
    Total         int    `json:"total"`
    Start         int    `json:"start"`
    Display       int    `json:"display"`
    Items         []Item `json:"items"`
}
