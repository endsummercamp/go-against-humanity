package game

type Deck struct {
	Name              			string
	Black_cards             	[]Card
	White_cards             	[]Card
	LastExtractedCard 			*Card
}

type DeckMetadata struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Official    bool   `json:"official"`
}

type DeckData struct {
	Black []BlackCard	`json:"black"`
	White []WhiteCard	`json:"white"`
	Metadata map[string]DeckMetadata	 `json:"metadata"`
}