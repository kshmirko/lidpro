const db = require('better-sqlite3')(process.env.DB_NAME, {});
db.pragma('journal_mode = WAL');

// const db = new sqlite3.Database(process.env.DB_NAME, err=>{
//   if(err){
//     return console.error(err.message)
//   }
//   console.log("Успешное подключение к БД "+process.env.DB_NAME)
// })

db.exec(
  `
  CREATE TABLE IF NOT EXISTS Experiment(
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
  `
)
db.exec(`
  CREATE TABLE IF NOT EXISTS Measurement(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    PROFILE_TIME TIMESTAMP NOT NULL DEFAULT(datetime('now')),
    PROFILE_STOP_TIME TIMESTAMP NOT NULL DEFAULT(datetime('now')),
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
        
  );
`)

module.exports = {db}