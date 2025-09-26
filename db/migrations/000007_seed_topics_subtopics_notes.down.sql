DELETE FROM smart_notes WHERE lesson_id = 1;
DELETE FROM subtopics WHERE topic_id IN (SELECT id FROM topics WHERE lesson_id = 1);
DELETE FROM topics WHERE lesson_id = 1;
