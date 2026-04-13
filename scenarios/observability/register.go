package observability

func Register() {
	registerAPICompat()
	registerContainerLogging()
	registerCRDStatus()
	registerPDB()
	registerTerminationPolicy()
}
