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

type MeasPlot struct {
	Alt float64
	Ch1 float64
	Ch2 float64
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

func CreateMeasurement(e Measurement) (int, error) {
	con := db.GetConnection()
	log.Println(con)
	tx := con.Db.MustBegin()
	var id int64
	qry := `
		INSERT INTO 
		Measurement(ID, PROFILE_TIME, REP_RATE, PROF_LEN, PROF_CNT, PROF_DATA_DAT, PROF_DATA_DAK, EXPERIMENT_ID)
		VALUES (NULL, ?,?,?,?,?,?,?)`

	time_1 := e.ProfileTime.Format("2006-01-02 15:04")
	res, err := tx.Exec(qry, time_1, e.RepRate, e.ProfLen, e.ProfCnt, e.ProfDataDat, e.ProfDataDak, e.ExperimentId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		id, _ = res.LastInsertId()
		tx.Commit()
	}

	db.CloseConnection()
	return int(id), nil
}

func GetMeasurementsByExperimentId(id int64) ([]Measurement, error) {
	con := db.GetConnection()
	tx := con.Db.MustBegin()

	rows, err := tx.Queryx(`
		SELECT ID, PROFILE_TIME, REP_RATE, PROF_LEN, PROF_CNT, PROF_DATA_DAT, PROF_DATA_DAK 
		FROM Measurement 
		WHERE EXPERIMENT_ID=?;
		`, id)
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
	db.CloseConnection()
	return ret, err
}
