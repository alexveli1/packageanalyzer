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

type PkgShort struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
}

type VerificationInfo struct {
	RequestArgs struct {
		Pkgset1 string `json:"pkgset1"`
		Pkgset2 string `json:"pkgset2"`
	} `json:"request_args"`
	Length   int `json:"length"`
	Packages []struct {
		Pkgset1  string   `json:"pkgset1"`
		Pkgset2  string   `json:"pkgset2"`
		Package1 PkgShort `json:"package1"`
		Package2 PkgShort `json:"package2"`
	} `json:"packages"`
}
type ComparePackage map[string]PkgShort
type Compare map[string][]string
type CompareBranch map[string]Compare
type Sources map[string][]string
type BranchSources map[string]Sources
type MethodBranches map[string]BranchSources
