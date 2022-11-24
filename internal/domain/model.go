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

type Comparepack struct {
	Name     string `json:"name"`
	Version1 string `json:"version1"`
	Release1 string `json:"release1"`
	Version2 string `json:"version2"`
	Release2 string `json:"release2"`
}
