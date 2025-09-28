package menuconfig

// PagedResponse is a generic response wrapper for paginated data.
type PagedResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"totalCount"`
}

type Grade struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type GradeUpsert struct {
	Name string `json:"name"`
}

type Subject struct {
	ID        int64  `json:"id"`
	GradeID   int64  `json:"gradeId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type SubjectUpsert struct {
	GradeID int64  `json:"gradeId"`
	Name    string `json:"name"`
}

type Lesson struct {
	ID        int64  `json:"id"`
	SubjectID int64  `json:"subjectId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type LessonUpsert struct {
	SubjectID int64  `json:"subjectId"`
	Name      string `json:"name"`
}

type Topic struct {
	ID        int64  `json:"id"`
	LessonID  int64  `json:"lessonId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type TopicUpsert struct {
	LessonID int64  `json:"lessonId"`
	Name     string `json:"name"`
}

type Subtopic struct {
	ID        int64  `json:"id"`
	TopicID   int64  `json:"topicId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type SubtopicUpsert struct {
	TopicID int64  `json:"topicId"`
	Name    string `json:"name"`
}

type Tutor struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"createdAt"`
}

type TutorUpsert struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type Year struct {
	ID        int64  `json:"id"`
	Value     int    `json:"value"`
	CreatedAt string `json:"createdAt"`
}

type YearUpsert struct {
	Value int `json:"value"`
}

type Tutorial struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	URL       *string `json:"url"`
	CreatedAt string  `json:"createdAt"`
}

type TutorialUpsert struct {
	Name string  `json:"name"`
	URL  *string `json:"url"`
}
