package testpackage

type TestHandler struct {
	Name       string            `json:"Name" json:"Name,omitempty"`
	FamilyName string            `json:"FamilyName" json:"FamilyName,omitempty"`
	Properties map[string]string `json:"Properties" json:"Properties,omitempty"`
}

func (TestHandler) GetMyEntity(str string, i int) string {
	return "get called!"
}

func (TestHandler) GetsMyEntities(str string) []string {
	return []string{"A", "B", "C"}
}

func (TestHandler) Post(str string, str2 string) string {
	return "post called!"
}

func (TestHandler) Delete(id int) string {
	return "del called!"
}

func (TestHandler) Patch(str string, float float64) string {
	return "patch called"
}
