CREATE TABLE lessons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subject_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
);

INSERT INTO lessons (subject_id, name, image_url)
VALUES 
(1, 'Parameter', '/images/lessons/parameter.png'),
(1, 'Cube', '/images/lessons/cube.png'),
(2, 'Grammar', '/images/lessons/grammar.png');