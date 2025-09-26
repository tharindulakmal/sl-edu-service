-- Geometry -> Triangles
INSERT INTO smart_notes (
  lesson_id, topic_id, subtopic_id, sub_topic_name,
  image_def_url, definition, theory, image_theory_url,
  example, image_example_url, is_default
) VALUES
(1, 1, 1, 'Triangles',
 'https://upload.wikimedia.org/wikipedia/commons/7/7e/Triangle_illustration.svg',
 'A triangle is a polygon with three edges and three vertices.',
 'The sum of the angles in a triangle is always 180 degrees.',
 'https://upload.wikimedia.org/wikipedia/commons/d/d0/Triangle_with_angles.svg',
 'Example: A right triangle with sides 3,4,5 satisfies Pythagoras theorem.',
 'https://upload.wikimedia.org/wikipedia/commons/4/4f/Triangle-right.svg',
 FALSE);

-- Geometry -> Squares
INSERT INTO smart_notes (
  lesson_id, topic_id, subtopic_id, sub_topic_name,
  image_def_url, definition, theory, image_theory_url,
  example, image_example_url, is_default
) VALUES
(1, 1, 2, 'Squares',
 'https://upload.wikimedia.org/wikipedia/commons/3/3c/Square_example.svg',
 'A square is a regular quadrilateral with four equal sides and angles.',
 'All angles in a square are 90 degrees.',
 'https://upload.wikimedia.org/wikipedia/commons/f/f3/Square_angles.svg',
 'Example: If each side = 4 units, perimeter = 16, area = 16.',
 'https://upload.wikimedia.org/wikipedia/commons/2/25/Square_example_calc.svg',
 FALSE);

-- Algebra -> Linear Equations
INSERT INTO smart_notes (
  lesson_id, topic_id, subtopic_id, sub_topic_name,
  image_def_url, definition, theory, image_theory_url,
  example, image_example_url, is_default
) VALUES
(1, 2, 3, 'Linear Equations',
 'https://upload.wikimedia.org/wikipedia/commons/2/29/Linear_function.svg',
 'An equation that makes a straight line when graphed.',
 'Form: ax + b = 0.',
 'https://upload.wikimedia.org/wikipedia/commons/5/5d/Linear_equation_graph.svg',
 'Example: 2x + 3 = 7 → x = 2.',
 'https://upload.wikimedia.org/wikipedia/commons/3/3e/Simple_linear_equation.svg',
 FALSE);

-- Algebra -> Quadratic Equations
INSERT INTO smart_notes (
  lesson_id, topic_id, subtopic_id, sub_topic_name,
  image_def_url, definition, theory, image_theory_url,
  example, image_example_url, is_default
) VALUES
(1, 2, 4, 'Quadratic Equations',
 'https://upload.wikimedia.org/wikipedia/commons/f/f7/Quadratic_function_parabola.svg',
 'An equation of degree 2: ax² + bx + c = 0.',
 'Its graph is a parabola.',
 'https://upload.wikimedia.org/wikipedia/commons/3/3f/Quadratic_parabola.svg',
 'Example: x² - 5x + 6 = 0 → roots are 2 and 3.',
 'https://upload.wikimedia.org/wikipedia/commons/b/b9/Quadratic_equation_roots.svg',
 FALSE);
