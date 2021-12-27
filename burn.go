package burn

import (
	"log"
	"os"
)

type Burn struct {
	inject.Injector
	handlers []Handler
	action   Handler
	logger   *log.Logger
}

func New() *Burn {
	m := &Burn{Injector: inject.New(), action: func() {}, logger: log.New(os.Stdout, "[burn] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	return m
}
