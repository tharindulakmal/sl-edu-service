-- Revert menu configuration table adjustments

-- Subtopics revert
ALTER TABLE subtopics
    DROP INDEX uq_subtopics_topic_name,
    DROP INDEX idx_subtopics_topic_id;

ALTER TABLE subtopics
    ADD COLUMN sort_order INT NOT NULL DEFAULT 0 AFTER name,
    MODIFY created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE subtopics
    CHANGE name sub_topic_name VARCHAR(255) NOT NULL;

ALTER TABLE subtopics
    ADD INDEX idx_subtopics_topic (topic_id, sort_order);

ALTER TABLE subtopics
    ADD CONSTRAINT fk_subtopics_topic FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE;

-- Topics revert
ALTER TABLE topics
    DROP INDEX uq_topics_lesson_name,
    DROP INDEX idx_topics_lesson_id;

ALTER TABLE topics
    ADD COLUMN main_topic_name VARCHAR(255) NOT NULL AFTER name,
    ADD COLUMN sort_order INT NOT NULL DEFAULT 0 AFTER main_topic_name,
    MODIFY created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE topics SET main_topic_name = name;

ALTER TABLE topics
    CHANGE name topic_name VARCHAR(255) NOT NULL;

ALTER TABLE topics
    ADD INDEX idx_topics_lesson (lesson_id, sort_order);

ALTER TABLE topics
    ADD CONSTRAINT fk_topics_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE;

-- Lessons revert
ALTER TABLE lessons
    DROP INDEX uq_lessons_subject_name,
    DROP INDEX idx_lessons_subject_id;

ALTER TABLE lessons
    DROP COLUMN created_at;

ALTER TABLE lessons
    ADD COLUMN image_url VARCHAR(255) NOT NULL AFTER name;

ALTER TABLE lessons
    ADD CONSTRAINT lessons_ibfk_1 FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE;

-- Subjects revert
ALTER TABLE subjects
    DROP INDEX uq_subjects_grade_name,
    DROP INDEX idx_subjects_grade_id;

ALTER TABLE subjects
    DROP COLUMN created_at;

ALTER TABLE subjects
    ADD CONSTRAINT subjects_ibfk_1 FOREIGN KEY (grade_id) REFERENCES grades(id) ON DELETE CASCADE;

-- Grades revert
ALTER TABLE grades
    DROP INDEX uq_grades_name;

ALTER TABLE grades
    DROP COLUMN created_at;

ALTER TABLE grades
    CHANGE name grade VARCHAR(50) NOT NULL;
