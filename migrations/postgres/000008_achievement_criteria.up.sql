CREATE TABLE IF NOT EXISTS achievement_criteria(
    achievemet_id INT REFERENCES achievements(id),
    stat_id INT REFERENCES stats(id),
    target_value INT
);