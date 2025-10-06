SET @subject_unique := (
  SELECT INDEX_NAME
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'subjects'
    AND INDEX_NAME = 'uq_subject_name'
  LIMIT 1
);
SET @drop_subject_unique := IF(@subject_unique IS NOT NULL,
  'ALTER TABLE subjects DROP INDEX uq_subject_name',
  'SELECT 1'
);
PREPARE drop_subject_stmt FROM @drop_subject_unique;
EXECUTE drop_subject_stmt;
DEALLOCATE PREPARE drop_subject_stmt;

SET @subject_grade_unique := (
  SELECT INDEX_NAME
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'subjects'
    AND INDEX_NAME = 'uq_subject_grade_name'
  LIMIT 1
);
SET @add_subject_grade_unique := IF(@subject_grade_unique IS NULL,
  'ALTER TABLE subjects ADD UNIQUE KEY uq_subject_grade_name (grade_id, name)',
  'SELECT 1'
);
PREPARE add_subject_grade_stmt FROM @add_subject_grade_unique;
EXECUTE add_subject_grade_stmt;
DEALLOCATE PREPARE add_subject_grade_stmt;

-- Drop junction
DROP TABLE IF EXISTS grade_subjects;
