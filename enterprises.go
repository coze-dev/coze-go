package coze

type enterprises struct {
	core    *core
	Members *enterprisesMembers
}

func newEnterprises(core *core) *enterprises {
	return &enterprises{
		core:    core,
		Members: newEnterprisesMembers(core),
	}
}
