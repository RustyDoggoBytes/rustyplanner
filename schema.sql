CREATE TABLE IF NOT EXISTS meals (
    user_id INTEGER NOT NULL,
    day DATE NOT NULL,
    breakfast text NOT NULL,
    snack1 text NOT NULL,
    lunch text NOT NULL,
    snack2 text NOT NULL,
    dinner text NOT NULL,
    PRIMARY KEY (user_id, day)
);
