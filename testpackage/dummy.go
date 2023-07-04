package testpackage

type TestHandler struct {
}

func (TestHandler) GetMyEntity(str string, i int) string {
	return "get called!"
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
