ALTER TABLE questions
    ADD COLUMN grade_id INT NULL AFTER lesson_id;

UPDATE questions q
LEFT JOIN lessons l ON q.lesson_id = l.id
LEFT JOIN subjects s ON l.subject_id = s.id
SET q.grade_id = s.grade_id
WHERE q.grade_id IS NULL;

ALTER TABLE questions
    MODIFY COLUMN grade_id INT NOT NULL,
    ADD INDEX idx_q_by_grade (grade_id, id);
