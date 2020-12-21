package source

import "log"

type Source interface {
	Sync() error
	Name() string
}

var monitor = make([]Source, 0)

func Monitoring(source Source) {
	monitor = append(monitor, source)
}

func Sync() {
	log.Printf("%d sources has sync\n", len(monitor))
	for _, source := range monitor {
		if source != nil {
			er := source.Sync()
			if er != nil {
				log.Printf("source : %s sync has err : %v\n", source.Name(), er)
			}
		}
	}
}
