CREATE TABLE IF NOT EXISTS grade_subjects (
  grade_id   INT NOT NULL,
  subject_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_grade_subject (grade_id, subject_id),
  KEY idx_gs_grade (grade_id),
  KEY idx_gs_subject (subject_id)
);

-- Normalise legacy subject names before deduplication
UPDATE subjects SET name = TRIM(name);

-- Cache the canonical subject row (lowest id) per distinct name
DROP TEMPORARY TABLE IF EXISTS tmp_subject_canonical;
CREATE TEMPORARY TABLE tmp_subject_canonical (
  name VARCHAR(120) NOT NULL PRIMARY KEY,
  canonical_id INT NOT NULL
) ENGINE = MEMORY;

INSERT INTO tmp_subject_canonical (name, canonical_id)
SELECT name, MIN(id) AS canonical_id
FROM subjects
GROUP BY name;

-- Populate the junction table using canonical subject ids
INSERT IGNORE INTO grade_subjects (grade_id, subject_id, created_at)
SELECT s.grade_id, t.canonical_id, COALESCE(s.created_at, NOW())
FROM subjects s
JOIN tmp_subject_canonical t ON s.name = t.name
WHERE s.grade_id IS NOT NULL;

-- Drop the lessons uniqueness constraint temporarily (it will be recreated later)
SET @lessons_unique := (
  SELECT INDEX_NAME
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'lessons'
    AND INDEX_NAME = 'uq_lessons_subject_name'
  LIMIT 1
);
SET @drop_lessons_unique := IF(@lessons_unique IS NOT NULL,
  'ALTER TABLE lessons DROP INDEX uq_lessons_subject_name',
  'SELECT 1'
);
PREPARE drop_lessons_unique_stmt FROM @drop_lessons_unique;
EXECUTE drop_lessons_unique_stmt;
DEALLOCATE PREPARE drop_lessons_unique_stmt;

-- Re-point lessons to the canonical subject id
UPDATE lessons l
JOIN subjects s ON l.subject_id = s.id
JOIN tmp_subject_canonical t ON s.name = t.name
SET l.subject_id = t.canonical_id;

-- If multiple lessons collapse to the same (subject_id, name) keep the oldest row
DELETE l
FROM lessons l
JOIN lessons keep_l
  ON l.subject_id = keep_l.subject_id
 AND l.name = keep_l.name
 AND l.id > keep_l.id;

-- Reinstate the lessons uniqueness constraint if it is absent
SET @lessons_unique_missing := (
  SELECT COUNT(*)
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'lessons'
    AND INDEX_NAME = 'uq_lessons_subject_name'
);
SET @create_lessons_unique := IF(@lessons_unique_missing = 0,
  'ALTER TABLE lessons ADD UNIQUE KEY uq_lessons_subject_name (subject_id, name)',
  'SELECT 1'
);
PREPARE create_lessons_unique_stmt FROM @create_lessons_unique;
EXECUTE create_lessons_unique_stmt;
DEALLOCATE PREPARE create_lessons_unique_stmt;

-- Remove duplicate subject rows now that dependants point at the canonical id
DELETE s
FROM subjects s
JOIN tmp_subject_canonical t ON s.name = t.name
WHERE s.id <> t.canonical_id;

-- Drop any legacy grade/name unique index before adding the new constraint
SET @legacy_subject_unique := (
  SELECT INDEX_NAME
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'subjects'
    AND INDEX_NAME IN ('uq_subject_grade_name', 'uq_subjects_grade_name')
  LIMIT 1
);
SET @drop_subject_grade_unique := IF(@legacy_subject_unique IS NOT NULL,
  CONCAT('ALTER TABLE subjects DROP INDEX ', @legacy_subject_unique),
  'SELECT 1'
);
PREPARE drop_subject_grade_unique_stmt FROM @drop_subject_grade_unique;
EXECUTE drop_subject_grade_unique_stmt;
DEALLOCATE PREPARE drop_subject_grade_unique_stmt;

-- Ensure the new unique constraint exists on subject name only
SET @subject_unique := (
  SELECT INDEX_NAME
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'subjects'
    AND INDEX_NAME = 'uq_subject_name'
  LIMIT 1
);
SET @add_subject_unique := IF(@subject_unique IS NULL,
  'ALTER TABLE subjects ADD UNIQUE KEY uq_subject_name (name)',
  'SELECT 1'
);
PREPARE add_subject_unique_stmt FROM @add_subject_unique;
EXECUTE add_subject_unique_stmt;
DEALLOCATE PREPARE add_subject_unique_stmt;

DROP TEMPORARY TABLE IF EXISTS tmp_subject_canonical;
