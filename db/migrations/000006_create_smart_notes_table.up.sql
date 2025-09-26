CREATE TABLE smart_notes (
  id INT AUTO_INCREMENT PRIMARY KEY,

  -- Scope of this smart note
  lesson_id INT NOT NULL,
  topic_id INT NULL,
  subtopic_id INT NULL,

  -- Render payload
  sub_topic_name VARCHAR(255) NOT NULL,
  image_def_url VARCHAR(512),
  definition TEXT,
  theory TEXT,
  image_theory_url VARCHAR(512),
  example TEXT,
  image_example_url VARCHAR(512),

  -- optional marker: default for a lesson (topic_id/subtopic_id may be NULL)
  is_default BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_notes_lookup (lesson_id, topic_id, subtopic_id),
  INDEX idx_notes_default (lesson_id, is_default),

  CONSTRAINT fk_notes_lesson   FOREIGN KEY (lesson_id)  REFERENCES lessons(id)   ON DELETE CASCADE,
  CONSTRAINT fk_notes_topic    FOREIGN KEY (topic_id)   REFERENCES topics(id)    ON DELETE SET NULL,
  CONSTRAINT fk_notes_subtopic FOREIGN KEY (subtopic_id)REFERENCES subtopics(id) ON DELETE SET NULL
);
