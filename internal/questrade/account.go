package questrade

type QuestradeAccount struct {
	Type   string `json:"type"`
	Number int64  `json:"number"`
	Status string
}

func GetAccounts() {}
