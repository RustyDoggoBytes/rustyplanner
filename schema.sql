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
    name VARCHAR(255) NOT NULL,
    completed BOOLEAN NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER IF NOT EXISTS update_groceries_last_updated
AFTER UPDATE ON groceries
FOR EACH ROW
    BEGIN
    UPDATE groceries SET last_updated = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

CREATE TABLE IF NOT EXISTS chores (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    recurrence_type VARCHAR(10) CHECK(recurrence_type IN ('once', 'weekly', 'monthly')) NOT NULL,
    recurrence_id INT NOT NULL,
    assigned VARCHAR(255),
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chores_recurrence_once(
    id INTEGER PRIMARY KEY,
    due_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS chores_recurrence_daily(
    id INTEGER PRIMARY KEY,
    frequency INT NOT NULL
);

CREATE TABLE IF NOT EXISTS chores_recurrence_weekly(
    id INTEGER PRIMARY KEY,
    frequency INT NOT NULL,
    mon INTEGER DEFAULT 0,
    tue INTEGER DEFAULT 0,
    wed INTEGER DEFAULT 0,
    thu INTEGER DEFAULT 0,
    fri INTEGER DEFAULT 0,
    sat INTEGER DEFAULT 0,
    sun INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS chores_recurrence_monthly(
    id INTEGER PRIMARY KEY,
    frequency INT NOT NULL,
    day_of_month INT NOT NULL CHECK (day_of_month BETWEEN 1 AND 31)
);


