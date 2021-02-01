package analyzer

import (
	"sync"

	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/util/config"

	"github.com/edersonbrilhante/ccvs/analyzer/anchoreengine"
	"github.com/edersonbrilhante/ccvs/analyzer/clair"
	"github.com/edersonbrilhante/ccvs/analyzer/trivy"
)

// StartAnalysis is function to execute analysis
func StartAnalysis(v *config.Vendors, al *ccvs.Analysis) {
	var wg sync.WaitGroup

	aa := []analyzer{
		&anchoreengine.Anchore{Config: v.AnchoreEngine},
		&clair.Clair{Config: v.Clair},
		&trivy.Trivy{Config: v.Trivy},
	}

	for _, a := range aa {
		wg.Add(1)
		go start(a, al, &wg)
	}

	wg.Wait()
}

type analyzer interface {
	Analyzer(*ccvs.Analysis) error
	Parser() error
}

func start(v analyzer, al *ccvs.Analysis, wg *sync.WaitGroup) {
	defer wg.Done()

	err := v.Analyzer(al)
	if err != nil {
		al.Errors = append(al.Errors, err.Error())
		return
	}
	err = v.Parser()
	if err != nil {
		al.Errors = append(al.Errors, err.Error())
	}
}
