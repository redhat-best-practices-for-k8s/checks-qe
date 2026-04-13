package platform

func Register() {
	registerHugepages()
	registerOCP()
	registerProbeChecks()
	registerServiceMesh()
}
