const {db} = require("./../db")


function getMeasurementsForExperimentId(id){

}

function insertMeasurements(meas){
    const stmt = db.prepare(`
    INSERT INTO 
        Measurement(PROFILE_TIME, REP_RATE, PROF_LEN, PROF_CNT, PROF_DATA_DAT, PROF_DATA_DAK, EXPERIMENT_ID)
    VALUES(?,?,?,?,?,?,?);
    `)

    //let res
    const insertMeas = db.transaction((meas)=>{
        let res = []
        for(const m of meas){
            const r = stmt.run(m.ProTime, m.RepRate, m.ProfLen, m.ProfCnt, m.ProfDataDat, m.ProfDataDak, m.ExpId)
            res.push(r)
        }
        return res
    })
    
    return insertMeas(rec)
} 


module.exports = {getMeasurementsForExperimentId, insertMeasurements}
