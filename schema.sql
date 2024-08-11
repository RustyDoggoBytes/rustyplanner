CREATE TABLE IF NOT EXISTS meals (
    user_id INTEGER NOT NULL,
    day DATE NOT NULL,
    breakfast TEXT NOT NULL,
    snack1 TEXT NOT NULL,
    lunch TEXT NOT NULL,
    snack2 TEXT NOT NULL,
    dinner TEXT NOT NULL,
    PRIMARY KEY (user_id, day)
);

CREATE TABLE IF NOT EXISTS groceries (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_groceries_last_updated
AFTER UPDATE ON groceries
FOR EACH ROW
    BEGIN
    UPDATE groceries SET last_updated = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;