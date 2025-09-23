CREATE TABLE subjects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  grade_id INT NOT NULL,
  FOREIGN KEY (grade_id) REFERENCES grades(id) ON DELETE CASCADE
);

-- Seed example data
INSERT INTO subjects (name, grade_id) VALUES
('Maths', 1),
('Sinhala', 1),
('Science', 2),
('English', 2);
