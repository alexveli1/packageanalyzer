package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
)

import (
	"context"
)

var sisyphus = []byte(`{
    "request_args": {
        "arch": null
    },
    "length": 179882,
    "packages": [{
            "name": "0ad-debuginfo",
            "epoch": 1,
            "version": "0.0.26",
            "release": "alt0_3_alpha",
            "arch": "aarch64",
            "disttag": "sisyphus+310089.100.1.1",
            "buildtime": 1668572449,
            "source": "0ad"
        }, {
            "name": "389-ds-base",
            "epoch": 0,
            "version": "2.2.3",
            "release": "alt2",
            "arch": "aarch64",
            "disttag": "sisyphus+308905.400.1.1",
            "buildtime": 1666675076,
            "source": "389-ds-base"
        }, {
            "name": "389-ds-base-debuginfo",
            "epoch": 0,
            "version": "2.2.3",
            "release": "alt2",
            "arch": "aarch64",
            "disttag": "sisyphus+308905.400.1.1",
            "buildtime": 1666675076,
            "source": "389-ds-base"
        }, {
            "name": "389-ds-base-devel",
            "epoch": 0,
            "version": "2.2.3",
            "release": "alt2",
            "arch": "aarch64",
            "disttag": "sisyphus+308905.400.1.1",
            "buildtime": 1666675076,
            "source": "389-ds-base"
        }, {
            "name": "jetty-http",
            "epoch": 0,
            "version": "9.4.40",
            "release": "alt1_2jpp11",
            "arch": "noarch",
            "disttag": "sisyphus+295303.100.1.1",
            "buildtime": 1644764356,
            "source": "jetty"
        }, {
            "name": "libXbae-docs",
            "epoch": 0,
            "version": "4.60.4",
            "release": "alt2.qa2",
            "arch": "noarch",
            "disttag": "",
            "buildtime": 1366177726,
            "source": "Xbae"
        }
    ]
}`)
var p10 = []byte(`{
			"request_args": {
        		"arch": null
			},
    		"length": 179882,
			"packages":[{
				"name": "0ad-debuginfo",
				"epoch": 1,
				"version": "0.0.26",
				"release": "alt0_3_alpha",
				"arch": "aarch64",
				"disttag": "sisyphus+310089.100.1.1",
				"buildtime": 1668572449,
				"source": "0ad"
			}, {
				"name": "389-ds-base",
				"epoch": 0,
				"version": "2.2.3",
				"release": "alt2",
				"arch": "aarch64",
				"disttag": "sisyphus+308905.400.1.1",
				"buildtime": 1666675076,
				"source": "389-ds-base"
			}, {
				"name": "389-ds-base-debuginfo",
				"epoch": 0,
				"version": "2.2.3",
				"release": "alt2",
				"arch": "aarch64",
				"disttag": "sisyphus+308905.400.1.1",
				"buildtime": 1666675076,
				"source": "389-ds-base"
			}, {
				"name": "389-ds-base-devel",
				"epoch": 0,
				"version": "2.2.3",
				"release": "alt2",
				"arch": "aarch64",
				"disttag": "sisyphus+308905.400.1.1",
				"buildtime": 1666675076,
				"source": "389-ds-base"
			}, {
				"name": "jetty-http",
				"epoch": 0,
				"version": "9.4.40",
				"release": "alt1_2jpp11",
				"arch": "noarch",
				"disttag": "sisyphus+295303.100.1.1",
				"buildtime": 1644764356,
				"source": "jetty"
			}, {
				"name": "libXbae-docs",
				"epoch": 0,
				"version": "4.60.4",
				"release": "alt2.qa2",
				"arch": "noarch",
				"disttag": "",
				"buildtime": 1366177726,
				"source": "Xbae"
			}]
		}`)

type ServiceTestSuite struct {
	suite.Suite
}

func TestAnalyzerService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestAnalyzerService_GetHigher() {
	type rpms struct {
		higher1 []byte
		higher2 []byte
		lower1  []byte
		lower2  []byte
	}
	tests := []struct {
		name    string
		branch1 string
		branch2 string
		rpms    rpms
		want    domain.Result
		wantErr bool
	}{
		{
			name:    "higher 0ad version",
			branch1: domain.Sisyphus,
			branch2: domain.P10,
			rpms: rpms{
				higher1: []byte(`[{"name": "0ad",
				"epoch": 1,
				"version": "0.0.27",
				"release": "alt0_3_alpha",
				"arch": "aarch64",
				"disttag": "sisyphus+310089.100.1.1",
				"buildtime": 1668572449,
				"source": "0ad"
			}]`), lower1: []byte(`[{
				"name": "0ad",
				"epoch": 1,
				"version": "0.0.26",
				"release": "alt0_3_alpha",
				"arch": "aarch64",
				"disttag": "sisyphus+310089.100.1.1",
				"buildtime": 1668572449,
				"source": "0ad"
			}]`),
				higher2: nil,
				lower2:  nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	repo := repository.NewRepositories()
	ctx := context.Background()
	for _, tt := range tests {
		s.Run(tt.name, func() {
			as := &AnalyzerService{
				repo: repo,
			}
			packSisyphus, err := prepBinpacks(sisyphus, tt.rpms.higher1)
			if err != nil {
				s.NoErrorf(err, "cannot prepare pack sisyphus")

				return
			}
			packP10, err := prepBinpacks(p10, tt.rpms.lower1)
			if err != nil {
				s.NoErrorf(err, "cannot prepare pack p10")

				return
			}
			err = as.repo.SavePacks(ctx, domain.Sisyphus, packSisyphus, nil)
			if err != nil {
				s.NoErrorf(err, "cannot save pack sisyphus")

				return
			}
			err = as.repo.SavePacks(ctx, domain.P10, packP10, nil)
			if err != nil {
				s.NoErrorf(err, "cannot save pack p10")

				return
			}
			var b []domain.Binpack
			err = json.Unmarshal(tt.rpms.higher1, &b)
			s.NoErrorf(err, "Error unmarshalling additional package: %v", err)
			got, err := as.GetHigher(ctx, tt.branch1, tt.branch2)
			s.NoErrorf(err, "unexpected error getting higher packages: %v", err)
			s.Equalf(got["aarch64"]["higher"]["sisyphus"], b, "unexpected resulting binpack set %v != %v", got["aarch64"]["higher"]["sisyphus"], b)
		})
	}
}
func (s *ServiceTestSuite) TestAnalyzerService_GetUnique() {
	type rpms struct {
		higher1 []byte
		higher2 []byte
		lower1  []byte
		lower2  []byte
		unique1 []byte
	}
	tests := []struct {
		name    string
		branch1 string
		branch2 string
		rpms    rpms
		want    domain.Result
		wantErr bool
	}{
		{
			name:    "unique 9wm version",
			branch1: domain.Sisyphus,
			branch2: domain.P10,
			rpms: rpms{
				unique1: []byte(`[{
            "name": "9wm",
            "epoch": 0,
            "version": "1.4.1",
            "release": "alt2",
            "arch": "aarch64",
            "disttag": "sisyphus+259420.100.1.1",
            "buildtime": 1602159269,
            "source": "9wm"
        }]`),
			},
			want:    nil,
			wantErr: false,
		},
	}
	repo := repository.NewRepositories()
	ctx := context.Background()
	for _, tt := range tests {
		s.Run(tt.name, func() {
			as := &AnalyzerService{
				repo: repo,
			}
			packSisyphus, err := prepBinpacks(sisyphus, tt.rpms.unique1)
			if err != nil {
				s.NoErrorf(err, "cannot prepare pack sisyphus")

				return
			}
			packP10, err := prepBinpacks(p10)
			if err != nil {
				s.NoErrorf(err, "cannot prepare pack p10")

				return
			}
			err = as.repo.SavePacks(ctx, domain.Sisyphus, packSisyphus, nil)
			if err != nil {
				s.NoErrorf(err, "cannot save pack sisyphus")

				return
			}
			err = as.repo.SavePacks(ctx, domain.P10, packP10, nil)
			if err != nil {
				s.NoErrorf(err, "cannot save pack p10")

				return
			}
			var b []domain.Binpack
			err = json.Unmarshal(tt.rpms.unique1, &b)
			s.NoErrorf(err, "Error unmarshalling additional package: %v", err)
			got, err := as.GetUnique(ctx, tt.branch1, tt.branch2)
			s.NoErrorf(err, "unexpected error getting higher packages: %v", err)
			check := got[b[0].Arch]["unique"]["sisyphus"]
			s.Equalf(check, b, "unexpected resulting binpack set %v != %v", check, b)
		})
	}
}

func prepBinpacks(branch []byte, data ...[]byte) (domain.Branch, error) {
	p := make([]domain.Binpack, 0)
	br := make(domain.Branch)
	var resp domain.RequestResult
	err := json.Unmarshal(branch, &resp)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(data); i++ {
		err = json.Unmarshal(data[i], &p)
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(p); j++ {
			resp.Packages = append(resp.Packages, p[j])
		}
	}
	for i := 0; i < len(resp.Packages); i++ {
		br[resp.Packages[i].Name] = append(br[resp.Packages[i].Name], resp.Packages[i])
	}

	return br, nil
}
