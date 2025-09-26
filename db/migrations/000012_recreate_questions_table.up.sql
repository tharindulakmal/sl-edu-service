DROP TABLE IF EXISTS questions;

CREATE TABLE questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lesson_id INT NOT NULL,
    topic_id INT NULL,
    subtopic_id INT NULL,
    tutor_id INT NULL,
    tute_id INT NULL,
    question TEXT NOT NULL,
    question_img_url VARCHAR(512),
    correct_answer VARCHAR(255) NOT NULL,
    theory TEXT,
    solution TEXT,
    other_answers JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- indexes for flexible queries and pagination
    INDEX idx_q_by_lesson   (lesson_id, id),
    INDEX idx_q_by_topic    (topic_id, id),
    INDEX idx_q_by_subtopic (subtopic_id, id),
    INDEX idx_q_by_tutor    (tutor_id, id),
    INDEX idx_q_by_tute     (tute_id, id)
);


INSERT INTO questions (lesson_id, question, question_img_url, correct_answer, theory, solution, other_answers)
VALUES
(1,
 'What is the perimeter of a square with side length 5 cm?',
 'https://upload.wikimedia.org/wikipedia/commons/3/3b/Square_-_black_simple.svg',
 '20 cm',
 'The perimeter of a square is calculated using the formula: P = 4 × side length.',
 'Here, side length = 5 cm. So, P = 4 × 5 = 20 cm.',
 JSON_ARRAY('10 cm', '15 cm', '20 cm', '25 cm')
),
(1,
 'If one side of a right triangle is 3 and the other is 4, what is the hypotenuse?',
 'https://upload.wikimedia.org/wikipedia/commons/9/9e/Right_triangle_with_notation_2.svg',
 '5',
 'By Pythagoras theorem: a² + b² = c².',
 '3² + 4² = 9 + 16 = 25. So, c = √25 = 5.',
 JSON_ARRAY('4', '5', '6', '7')
),
(1,
 'Which of the following is NOT a prime number?',
 NULL,
 '9',
 'Prime numbers are natural numbers greater than 1 with only two factors: 1 and itself.',
 '9 is not prime because it can be divided evenly by 3. (9 ÷ 3 = 3).',
 JSON_ARRAY('2', '3', '9', '11')
),
(1,
 'The area of a rectangle with length 8 cm and width 6 cm is?',
 'https://upload.wikimedia.org/wikipedia/commons/4/47/Rectangle_example.svg',
 '48 cm²',
 'The area of a rectangle is given by A = length × width.',
 'A = 8 × 6 = 48 cm².',
 JSON_ARRAY('14 cm²', '24 cm²', '30 cm²', '48 cm²')
);