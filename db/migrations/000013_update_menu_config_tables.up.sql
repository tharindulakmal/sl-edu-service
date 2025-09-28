-- Align legacy menu configuration tables with new repository expectations

-- Grades adjustments
ALTER TABLE grades
    CHANGE grade name VARCHAR(120) NOT NULL;

ALTER TABLE grades
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER name;

ALTER TABLE grades
    ADD UNIQUE INDEX uq_grades_name (name);

-- Subjects adjustments
ALTER TABLE subjects
    DROP FOREIGN KEY subjects_ibfk_1;

ALTER TABLE subjects
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER name;

CREATE UNIQUE INDEX uq_subjects_grade_name ON subjects (grade_id, name);
CREATE INDEX idx_subjects_grade_id ON subjects (grade_id);

-- Lessons adjustments
ALTER TABLE lessons
    DROP FOREIGN KEY lessons_ibfk_1;

ALTER TABLE lessons
    DROP COLUMN image_url;

ALTER TABLE lessons
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER name;

CREATE UNIQUE INDEX uq_lessons_subject_name ON lessons (subject_id, name);
CREATE INDEX idx_lessons_subject_id ON lessons (subject_id);

-- Topics adjustments
ALTER TABLE topics
    DROP FOREIGN KEY fk_topics_lesson;

ALTER TABLE topics
    DROP INDEX idx_topics_lesson;

ALTER TABLE topics
    CHANGE topic_name name VARCHAR(255) NOT NULL;

ALTER TABLE topics
    DROP COLUMN main_topic_name,
    DROP COLUMN sort_order,
    MODIFY created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE UNIQUE INDEX uq_topics_lesson_name ON topics (lesson_id, name);
CREATE INDEX idx_topics_lesson_id ON topics (lesson_id);

-- Subtopics adjustments
ALTER TABLE subtopics
    DROP FOREIGN KEY fk_subtopics_topic;

ALTER TABLE subtopics
    DROP INDEX idx_subtopics_topic;

ALTER TABLE subtopics
    CHANGE sub_topic_name name VARCHAR(255) NOT NULL;

ALTER TABLE subtopics
    DROP COLUMN sort_order,
    MODIFY created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE UNIQUE INDEX uq_subtopics_topic_name ON subtopics (topic_id, name);
CREATE INDEX idx_subtopics_topic_id ON subtopics (topic_id);
