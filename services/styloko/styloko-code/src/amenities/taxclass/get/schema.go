package get

type TaxClassResponse struct {
	Id         int     `json:"taxClassId"`
	Name       string  `json:"name"`
	IsDefault  bool    `json:"isDefault"`
	TaxPercent float64 `json:"taxPercent"`
}
