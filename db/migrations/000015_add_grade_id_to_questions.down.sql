ALTER TABLE questions
    DROP INDEX idx_q_by_grade,
    DROP COLUMN grade_id;
