CREATE TABLE subtopics (
  id INT AUTO_INCREMENT PRIMARY KEY,
  topic_id INT NOT NULL,
  sub_topic_name VARCHAR(255) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_subtopics_topic (topic_id, sort_order),
  CONSTRAINT fk_subtopics_topic FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
);
