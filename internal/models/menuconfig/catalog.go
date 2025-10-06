package menuconfig

type CatalogResponse struct {
	Grades        []Grade        `json:"grades"`
	Subjects      []Subject      `json:"subjects"`
	GradeSubjects []GradeSubject `json:"gradeSubjects"`
	Lessons       []Lesson       `json:"lessons"`
	Topics        []Topic        `json:"topics"`
	Subtopics     []Subtopic     `json:"subtopics"`
}

type GradeSubject struct {
	GradeID   int64  `json:"gradeId"`
	SubjectID int64  `json:"subjectId"`
	CreatedAt string `json:"createdAt"`
}
