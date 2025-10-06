-- Create junction table to map Grades ↔ Subjects (no FKs for loose coupling)
CREATE TABLE IF NOT EXISTS grade_subjects (
  grade_id   INT NOT NULL,
  subject_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_grade_subject (grade_id, subject_id),
  KEY idx_gs_grade (grade_id),
  KEY idx_gs_subject (subject_id)
);

-- Ensure subjects have a global unique name (prevent “English” duplicates per grade)
-- If there is an existing unique on (grade_id, name), drop it first (ignore if absent)
ALTER TABLE subjects DROP INDEX uq_subject_grade_name;
ALTER TABLE subjects ADD UNIQUE KEY uq_subject_name (name);

-- Optional: if subjects.grade_id currently exists, keep it for compatibility now (do NOT drop here)

-- Backfill grade_subjects from existing subjects rows that still carry grade_id
INSERT IGNORE INTO grade_subjects (grade_id, subject_id, created_at)
SELECT DISTINCT s.grade_id, s.id, COALESCE(s.created_at, NOW())
FROM subjects s
WHERE s.grade_id IS NOT NULL;
