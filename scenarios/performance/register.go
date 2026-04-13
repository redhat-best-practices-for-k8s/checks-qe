package performance

func Register() {
	registerExecProbes()
	registerResources()
	registerSchedulingPolicy()
}
