package accesscontrol

func Register() {
	registerAutomountToken()
	registerCapabilities()
	registerCRDRoles()
	registerHostIPC()
	registerHostNetwork()
	registerHostPath()
	registerHostPID()
	registerHostPort()
	registerNamespace()
	registerProcesses()
	registerRBAC()
	registerRequests()
	registerResourceQuota()
	registerSecurityContext()
	registerServices()
	registerSysNice()
	registerSysPtrace()
}
