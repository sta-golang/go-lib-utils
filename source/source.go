package source

import "github.com/sta-golang/go-lib-utils/log"

type Source interface {
	Sync() error
	Name() string
}

var monitor = make([]Source, 0)

func Monitoring(source Source) {
	monitor = append(monitor, source)
}

func Sync() {
	log.FrameworkLogger.Infof("%d sources has sync", len(monitor))
	for _, source := range monitor {
		if source != nil {
			er := source.Sync()
			if er != nil {
				log.FrameworkLogger.Warnf("source : %s sync has err : %v", source.Name(), source.Sync())
			}
		}
	}
}
