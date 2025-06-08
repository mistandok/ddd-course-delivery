package cmd

import "log"

type Closer interface {
	Close() error
}

func (cr *CompositionRoot) RegisterCloser(c Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, closer := range cr.closers {
		if err := closer.Close(); err != nil {
			log.Printf("error closing resource: %v", err)
		}
	}
}
