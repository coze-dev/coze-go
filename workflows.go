package coze

type workflows struct {
	Runs *workflowRuns
	Chat *workflowsChat
}

func newWorkflows(core *core) *workflows {
	return &workflows{
		Runs: newWorkflowRun(core),
		Chat: newWorkflowsChat(core),
	}
}
