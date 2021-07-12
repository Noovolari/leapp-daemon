package providers

type Providers struct {
}

func NewProviders() *Providers {
	return &Providers{}
}

func (prov *Providers) Close() {
	prov.GetTimerCollection().Close()
}
