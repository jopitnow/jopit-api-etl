package dto

type EligibleDTO struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	IsRequired bool        `json:"is_required"`
	Options    []OptionDTO `json:"options"`
}

type Attributes map[string]string
type OptionDTO string
