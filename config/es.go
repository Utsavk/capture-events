package config

import (
	es "github.com/Utsavk/capture-events/db/es"
	"github.com/mitchellh/mapstructure"
)

var ESConf *es.Conf

func parseESConf() error {
	if conf, ok := Props.DB["ES"]; ok {
		ESConf = &es.Conf{}
		mapstructure.Decode(conf, ESConf)
	}
	return ESConf.ValidateConf()
}
