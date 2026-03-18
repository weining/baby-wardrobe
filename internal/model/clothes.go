package model

import "time"

// 类别
const (
	CategoryTop     = "上衣"
	CategoryBottom  = "裤子"
	CategoryOuter   = "外套"
	CategoryShoes   = "鞋子"
	CategoryPajamas = "睡衣"
	CategoryAccess  = "配件"
	CategoryOther   = "其他"
)

// 季节
const (
	SeasonSpringAutumn = "春秋"
	SeasonSummer       = "夏季"
	SeasonWinter       = "冬季"
	SeasonAllYear      = "四季"
)

// 状态
const (
	StatusWearing  = "wearing"
	StatusStored   = "stored"
	StatusOutgrown = "outgrown"
	StatusDonated  = "donated"
)

// Seasons 固定不变
var Seasons = []string{
	SeasonSpringAutumn, SeasonSummer, SeasonWinter, SeasonAllYear,
}

// DefaultCategories 数据库为空时的种子数据
var DefaultCategories = []string{
	CategoryTop, CategoryBottom, CategoryOuter,
	CategoryShoes, CategoryPajamas, CategoryAccess, CategoryOther,
}

// Status 表示一个衣物状态，存储在数据库中
type Status struct {
	ID    int64
	Value string // 唯一键，如 "wearing"
	Label string // 显示名，如 "在穿"
	Color string // green / blue / yellow / gray / red / purple
}

// DefaultStatuses 数据库为空时的种子数据
var DefaultStatuses = []Status{
	{Value: StatusWearing, Label: "在穿", Color: "green"},
	{Value: StatusStored, Label: "收纳中", Color: "blue"},
	{Value: StatusOutgrown, Label: "已小", Color: "yellow"},
	{Value: StatusDonated, Label: "已捐出", Color: "gray"},
}

type Cloth struct {
	ID        int64
	Name      string
	Category  string
	Size      string
	Season    string
	Status    string
	Color     string
	Brand     string
	PhotoPath string
	ThumbPath string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Cloth) StatusLabel() string { return c.Status }
func (c *Cloth) StatusColor() string { return "gray" }
