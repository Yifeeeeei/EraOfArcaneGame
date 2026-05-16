package model

// Card represents a card definition exposed to the game engine and frontend.
// Server-side definitions are implemented in Go; JSON snapshots are generated
// or reference data, not the runtime source of truth.
type Card struct {
	Number          string         `json:"number"`
	Type            string         `json:"type"` // 人物/伙伴/技能/道具
	Name            string         `json:"name"`
	Category        string         `json:"category"` // element category
	Tag             string         `json:"tag"`      // e.g. "装备-饰物"
	Description     string         `json:"description"`
	Quote           string         `json:"quote"`
	ElementsCost    map[string]int `json:"elements_cost"` // cost to play
	ElementsGain    map[string]int `json:"elements_gain"` // elements provided (负载)
	ElementsExpense map[string]int `json:"elements_expense"`
	VersionNum      string         `json:"version_number"`
	VersionName     string         `json:"version_name"`
	Attack          int            `json:"attack"`   // -1 = N/A
	Life            int            `json:"life"`     // -1 = N/A
	Duration        int            `json:"duration"` // -1 = N/A
	Power           int            `json:"power"`    // -1 = N/A
	Spawns          []string       `json:"spawns"`
	OutputPath      string         `json:"output_path"`
}

// CardType constants
const (
	CardTypeHero      = "人物"
	CardTypeCompanion = "伙伴"
	CardTypeSkill     = "技能"
	CardTypeItem      = "道具"
)

// Element constants (matching JSON data short names)
const (
	ElementFire   = "火"
	ElementWater  = "水"
	ElementEarth  = "地"
	ElementAir    = "气"
	ElementLight  = "光"
	ElementShadow = "暗"
	ElementArcane = "无" // 奥术/无 = wildcard element
)

var AllElements = []string{
	ElementFire, ElementWater, ElementEarth, ElementAir,
	ElementLight, ElementShadow, ElementArcane,
}

// ElementDisplayNames maps short element names to display names
var ElementDisplayNames = map[string]string{
	"火": "火焰", "水": "水纹", "地": "地脉", "气": "大气",
	"光": "光辉", "暗": "暗影", "无": "奥术",
}

// IsHero returns true if the card is a hero card
func (c *Card) IsHero() bool { return c.Type == CardTypeHero }

// IsCompanion returns true if the card is a companion card
func (c *Card) IsCompanion() bool { return c.Type == CardTypeCompanion }

// IsSkill returns true if the card is a skill card
func (c *Card) IsSkill() bool { return c.Type == CardTypeSkill }

// IsItem returns true if the card is an item card
func (c *Card) IsItem() bool { return c.Type == CardTypeItem }

// TotalCost returns the total element cost to play this card
func (c *Card) TotalCost() int {
	total := 0
	for _, v := range c.ElementsCost {
		total += v
	}
	return total
}

// IsTerrain checks if this item card is a terrain card (地形)
// Terrain cards are placed on the battlefield grid instead of being consumed
func (c *Card) IsTerrain() bool {
	return false
}
