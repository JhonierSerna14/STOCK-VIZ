package models

type StockResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}
