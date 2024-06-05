-- +goose Up
-- +goose StatementBegin
CREATE TABLE Measurement(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    PROFILE_TIME TIMESTAMP NOT NULL DEFAULT(datetime('now')),
    REP_RATE INTEGER NOT NULL DEFAULT(10),
    PROF_LEN INTEGER UNSIGNED NOT NULL DEFAULT(512),
    PROF_CNT INTEGER UNSIGNED NOT NULL DEFAULT(1),
    PROF_DATA_DAT TEXT NOT NULL,
    PROF_DATA_DAK TEXT NOT NULL,
    EXPERIMENT_ID INTEGER NOT NULL,
    CONSTRAINT
        chk_reprate_constraint CHECK((REP_RATE>=1) AND (REP_RATE<=100)),
    CONSTRAINT
        chk_proflen_constraint CHECK((PROF_LEN>=128) AND (PROF_LEN<=1024)),
    CONSTRAINT
        chk_profcnt_constraint CHECK((PROF_CNT>=1) AND (PROF_CNT<=18000)),
    CONSTRAINT
        chk_prof_data_dat_constraint CHECK(JSON_VALID(PROF_DATA_DAT)),
    CONSTRAINT
        chk_prof_data_dak_constraint CHECK(JSON_VALID(PROF_DATA_DAK)),
    CONSTRAINT
        fk_experiment_id FOREIGN KEY(EXPERIMENT_ID) REFERENCES Experiment(ID) DEFERRABLE INITIALLY DEFERRED
        
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Measurement;
-- +goose StatementEnd
