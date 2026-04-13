package builder

type MockCertValidator struct {
	ContainerCertified bool
	OperatorCertified  bool
	HelmChartCertified bool
}

func (m *MockCertValidator) IsContainerCertified(_, _, _, _ string) bool {
	return m.ContainerCertified
}

func (m *MockCertValidator) IsOperatorCertified(_, _ string) bool {
	return m.OperatorCertified
}

func (m *MockCertValidator) IsHelmChartCertified(_, _, _ string) bool {
	return m.HelmChartCertified
}
