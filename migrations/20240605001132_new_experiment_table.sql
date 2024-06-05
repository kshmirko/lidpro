-- +goose Up
-- +goose StatementBegin
CREATE TABLE Experiment(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    START_TIME TIMESTAMP DEFAULT(datetime('now')),
    TITLE VARCHAR(100) NOT NULL,
    COMMENT VARCHAR(500) NOT NULL,
    VERT_RES FLOAT NOT NULL DEFAULT(1500.0),
    ACCUM_TIME INTEGER UNSIGNED NOT NULL DEFAULT(10),
    ARCHIVE BLOB NOT NULL,
    CONSTRAINT
    chk_statial_step_constraint CHECK ((VERT_RES>=1500.0) AND (VERT_RES<=1912.5)),
    CONSTRAINT
    chk_accum_time_constraint CHECK ((ACCUM_TIME>=10) AND (ACCUM_TIME<=18000))
);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Experiment;
SELECT 'down SQL query';
-- +goose StatementEnd
