package models

type Student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

type Question struct {
	Id       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	TestId   string `json:"test_id"`
}

type Test struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Enrollment struct {
	StudentId string `json:"student_id"`
	TestId    string `json:"test_id"`
}
