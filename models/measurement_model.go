package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kshmirko/lidpro/db"
)

type Measurement struct {
	Id           int       `db:"ID"`
	ProfileTime  time.Time `db:"PROFILE_TIME"`
	RepRate      uint32    `db:"REP_RATE"`
	ProfLen      uint32    `db:"PROF_LEN"`
	ProfCnt      uint32    `db:"PROF_CNT"`
	ProfDataDat  string    `db:"PROF_DATA_DAT"`
	ProfDataDak  string    `db:"PROF_DATA_DAK"`
	ProfDat      []float64 // не отображаемые поля на SQL запрос
	ProfDak      []float64 // не отображаемые поля на SQL запрос
	ExperimentId int       `db:"EXPERIMENT_ID"`
}

func GetAllMeasurements() []Measurement {
	con := db.GetConnection()
	tx := con.Db.MustBegin()

	rows, err := tx.Queryx("SELECT * FROM Measurement")
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]Measurement, 0, 10)
	for rows.Next() {
		var e Measurement
		err = rows.StructScan(&e)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(e.ProfDataDat), &e.ProfDat)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(e.ProfDataDak), &e.ProfDak)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, e)
	}

	tx.Commit()
	db.CloseConnection()
	return ret
}
