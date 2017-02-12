package main

type Party struct {
	ID        string `json:"ID"`
	LeadID    string `json:"lead,omitempty"`
	Name      string `json:"name,omitempty"`
	NumPeople int    `json:"numPeople,omitempty'`
	Replied   bool   `json:"replied"`
	NumComing int    `json:"numComing,omitempty"`
	Address   string `json:"address,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	MagicWord string `json:"codeWord,omitempty"`
	SortValue string `json:"sortValue"`
}

type Person struct {
	ID                   string `json:"ID"`
	PartyID              string `json:"party,omitempty"`
	Name                 string `json:"name,omitempty"`
	GetsPlusOne          bool   `json:"getsPlusOne"`
	PlusOneID            string `json:"plusOne,omitempty"`
	IsPlusOne            bool   `json:"isPlusOne"`
	IsPlusOneOfID        string `json:"isPlusOneOf,omitempty"`
	Replied              bool   `json:"replied"`
	Reply                bool   `json:"reply"`
	DietaryRestrictions  string `json:"dietaryRestrictions,omitempty"`
	IsChild              bool   `json:"isChild"`
	WillAccompanyID      string `json:"willAccompany,omitempty"`
	BabysitterForWedding bool   `json:"babysitterForWedding"`
	BabysitterForEvents  bool   `json:"babysitterForEvents"`
}
