package coze

type stores struct {
	core    *core
	Plugins *storesPlugins
}

func newStores(core *core) *stores {
	return &stores{core: core, Plugins: newStoresPlugins(core)}
}
