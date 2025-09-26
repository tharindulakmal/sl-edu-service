-- assumes you already have lessons seeded; e.g., lesson id 1 = "Maths"

-- Topics for lesson 1 (Maths)
INSERT INTO topics (lesson_id, topic_name, main_topic_name, sort_order) VALUES
(1, 'Geometry',      'Introduction to Geometry', 1),
(1, 'Algebra',       'Equations',                 2),
(1, 'Trigonometry',  'Angles and Ratios',        3);

-- Subtopics
-- Geometry = topic_id 1 (check your actual IDs after insert; during dev this is fine)
INSERT INTO subtopics (topic_id, sub_topic_name, sort_order) VALUES
(1, 'Triangles', 1),
(1, 'Squares',   2);

-- Algebra = topic_id 2
INSERT INTO subtopics (topic_id, sub_topic_name, sort_order) VALUES
(2, 'Linear Equations',    1),
(2, 'Quadratic Equations', 2);

-- Default smart note for lesson 1 (Maths)
INSERT INTO smart_notes (
  lesson_id, topic_id, subtopic_id, sub_topic_name,
  image_def_url, definition, theory, image_theory_url, example, image_example_url, is_default
) VALUES
(1, NULL, NULL, 'Pythagoras Theorem',
 'https://upload.wikimedia.org/wikipedia/commons/d/d2/Pythagorean.svg',
 'In a right-angled triangle, the square of the hypotenuse is equal to the sum of the squares of the other two sides.',
 'The theorem states: a² + b² = c², where c is the hypotenuse of a right-angled triangle.',
 'https://upload.wikimedia.org/wikipedia/commons/2/2e/Pythagoras-2a.svg',
 'If sides are 3 and 4, hypotenuse is 5 since 3²+4²=25=5².',
 'https://upload.wikimedia.org/wikipedia/commons/9/9e/Right_triangle_with_notation_2.svg',
 TRUE
);
