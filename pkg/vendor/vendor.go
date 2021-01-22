package vendor

import (
	"sync"

	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/util/config"
)

// StartAnalysis is function to execute analysis
func StartAnalysis(v *config.Vendors, al *ccvs.Analysis) {
	var wg sync.WaitGroup

	aa := []vendor{
		&Anchore{config: v.AnchoreEngine},
	}

	for _, a := range aa {
		wg.Add(1)
		go start(a, al, &wg)
	}

	wg.Wait()
}

type vendor interface {
	Analyzer(*ccvs.Analysis) error
	Parser() error
}

func start(v vendor, al *ccvs.Analysis, wg *sync.WaitGroup) {
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
