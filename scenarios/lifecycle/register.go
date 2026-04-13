package lifecycle

func Register() {
	registerAffinity()
	registerHooks()
	registerPodConfig()
	registerPodRecreation()
	registerProbes()
	registerScaling()
	registerScheduling()
	registerStorage()
}
