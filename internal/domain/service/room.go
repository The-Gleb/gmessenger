package service

type Room struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Personal bool
	Clients  map[string]*Client `json:"clients"`
}
