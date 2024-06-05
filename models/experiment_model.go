package models

import (
	"log"
	"time"

	"github.com/kshmirko/lidpro/db"
)

type Experiment struct {
	Id        int       `db:"ID"`
	Title     string    `db:"TITLE"`
	Comment   string    `db:"COMMENT"`
	StartTime time.Time `db:"START_TIME"`
	VertRes   float32   `db:"VERT_RES"`
	AccumTime uint32    `db:"ACCUM_TIME"`
	Archive   []byte    `db:"ARCHIVE"`
}

func GetAllExperimentsWithoutArchive() []Experiment {
	con := db.GetConnection()
	tx := con.Db.MustBegin()

	rows, err := tx.Queryx("SELECT ID, TITLE, COMMENT, START_TIME, VERT_RES, ACCUM_TIME FROM Experiment")
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]Experiment, 0, 10)
	for rows.Next() {
		var e Experiment
		err = rows.StructScan(&e)
		if err != nil {
			log.Fatal(err)
		}

		ret = append(ret, e)
	}

	tx.Commit()
	con.Db.Close()
	return ret
}

func GetAllExperiments() []Experiment {
	con := db.GetConnection()
	tx := con.Db.MustBegin()

	rows, err := tx.Queryx("SELECT * FROM Experiment")
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]Experiment, 0, 10)
	for rows.Next() {
		var e Experiment
		err = rows.StructScan(&e)
		if err != nil {
			log.Fatal(err)
		}

		ret = append(ret, e)
	}

	tx.Commit()
	con.Db.Close()
	return ret
}
