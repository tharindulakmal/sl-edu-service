CREATE TABLE IF NOT EXISTS grades (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_grades_name (name)
);

CREATE TABLE IF NOT EXISTS subjects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    grade_id INT NOT NULL,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_subjects_grade_name (grade_id, name),
    INDEX idx_subjects_grade_id (grade_id)
);

CREATE TABLE IF NOT EXISTS lessons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subject_id INT NOT NULL,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_lessons_subject_name (subject_id, name),
    INDEX idx_lessons_subject_id (subject_id)
);

CREATE TABLE IF NOT EXISTS topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lesson_id INT NOT NULL,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_topics_lesson_name (lesson_id, name),
    INDEX idx_topics_lesson_id (lesson_id)
);

CREATE TABLE IF NOT EXISTS subtopics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    topic_id INT NOT NULL,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_subtopics_topic_name (topic_id, name),
    INDEX idx_subtopics_topic_id (topic_id)
);

CREATE TABLE IF NOT EXISTS tutors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    email VARCHAR(160),
    phone VARCHAR(40),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS years (
    id INT AUTO_INCREMENT PRIMARY KEY,
    value INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_years_value (value)
);

CREATE TABLE IF NOT EXISTS tutorials (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(160) NOT NULL,
    url VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
