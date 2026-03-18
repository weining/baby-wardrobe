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

var Categories = []string{
	CategoryTop, CategoryBottom, CategoryOuter,
	CategoryShoes, CategoryPajamas, CategoryAccess, CategoryOther,
}

var Seasons = []string{
	SeasonSpringAutumn, SeasonSummer, SeasonWinter, SeasonAllYear,
}

var StatusOptions = []struct {
	Value string
	Label string
	Color string
}{
	{StatusWearing, "在穿", "green"},
	{StatusStored, "收纳中", "blue"},
	{StatusOutgrown, "已小", "yellow"},
	{StatusDonated, "已捐出", "gray"},
}

func StatusLabel(status string) string {
	for _, s := range StatusOptions {
		if s.Value == status {
			return s.Label
		}
	}
	return status
}

func StatusColor(status string) string {
	for _, s := range StatusOptions {
		if s.Value == status {
			return s.Color
		}
	}
	return "gray"
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

func (c *Cloth) StatusLabel() string { return StatusLabel(c.Status) }
func (c *Cloth) StatusColor() string { return StatusColor(c.Status) }
