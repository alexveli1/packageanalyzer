package domain

type Binpack struct {
	Name      string `json:"name"`
	Epoch     int    `json:"epoch"`
	Version   string `json:"version"`
	Release   string `json:"release"`
	Arch      string `json:"arch"`
	Disttag   string `json:"disttag"`
	Buildtime int    `json:"buildtime"`
	Source    string `json:"source"`
}

type RequestResult struct {
	Request_args interface{} `json:"request_args"`
	Length       int64       `json:"length"`
	Packages     []Binpack   `json:"packages"`
}

type Result map[string]map[string]map[string][]Binpack
type Method map[string]map[string][]Binpack
type Branch map[string][]Binpack
