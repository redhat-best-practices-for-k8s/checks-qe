package networking

func Register() {
	registerICMP()
	registerNetworkServices()
	registerPolicies()
	registerPorts()
	registerSRIOV()
}
