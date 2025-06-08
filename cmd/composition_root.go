package cmd

type CompositionRoot struct {
	configs Config

	closers []Closer
}

func NewCompositionRoot(configs Config) CompositionRoot {
	return CompositionRoot{
		configs: configs,
	}
}
