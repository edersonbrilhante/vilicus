package analyzer

import (
	"sync"

	"github.com/edersonbrilhante/vilicus"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"

	"github.com/edersonbrilhante/vilicus/pkg/analyzer/anchoreengine"
	"github.com/edersonbrilhante/vilicus/pkg/analyzer/clair"
	"github.com/edersonbrilhante/vilicus/pkg/analyzer/trivy"
)

// StartAnalysis is function to execute analysis
func StartAnalysis(v *config.Vendors, al *vilicus.Analysis) {
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
	Analyzer(*vilicus.Analysis) error
	Parser() error
}

func start(v analyzer, al *vilicus.Analysis, wg *sync.WaitGroup) {
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
