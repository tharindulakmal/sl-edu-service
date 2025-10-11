CREATE TABLE IF NOT EXISTS grade_subject_lessons (
  grade_id   INT NOT NULL,
  subject_id INT NOT NULL,
  lesson_id  INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (grade_id, subject_id, lesson_id),
  KEY idx_gsl_grade_subject (grade_id, subject_id),
  KEY idx_gsl_lesson (lesson_id)
);
