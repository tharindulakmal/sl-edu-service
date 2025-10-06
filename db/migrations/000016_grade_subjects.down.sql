-- Revert unique back to (grade_id, name) if you had it previously (optional; best-effort)
ALTER TABLE subjects DROP INDEX uq_subject_name;
ALTER TABLE subjects ADD UNIQUE KEY uq_subject_grade_name (grade_id, name);

-- Drop junction
DROP TABLE IF EXISTS grade_subjects;
