const {db} = require("./../db")


/**
 * 
 * @param {number} id - идентификатор эксперимента 
 * @returns {List<Object>} список измерений для текущего эксперимента
 */
function getMeasurementsForExperimentId(id){
    const stmt = db.prepare(`
        SELECT 
            ID, PROFILE_TIME, PROFILE_STOP_TIME, REP_RATE, PROF_LEN, PROF_CNT, PROF_DATA_DAT, PROF_DATA_DAK
        FROM
            Measurement
        WHERE
            EXPERIMENT_ID = ?
        ORDER BY
            PROFILE_TIME ASC;
    `);

    let measurements = []
   
    for(const m of stmt.iterate(id)){
        const dat = m.PROF_DATA_DAT
        const dak = m.PROF_DATA_DAK
        measurements.push({
            Id: m.ID,
            ProfTime: new Date(m.PROFILE_TIME),
            StopTime: new Date(m.PROFILE_STOP_TIME),
            RepRate: m.REP_RATE,
            ProfLen: m.PROF_LEN,
            ProfCnt: m.PROF_CNT,
            Dat: JSON.parse(dat),
            Dak: JSON.parse(dak)
        })
    }

    return measurements
}

/**
 * 
 * @param {Object} meas - структура, опрееляющая измерение 
 * @param {number} expid - идентификатор сопуствующего эксперимента
 * @returns {List<Object>} - список вставленных объектов
 */
function insertMeasurements(meas, expid){
    const stmt = db.prepare(`
    INSERT INTO 
        Measurement(PROFILE_TIME, PROFILE_STOP_TIME, REP_RATE, PROF_LEN, PROF_CNT, PROF_DATA_DAT, PROF_DATA_DAK, EXPERIMENT_ID)
    VALUES
        (?,?,?,?,?,?,?,?);
    `)

 
    //let res
    const insertMeas = db.transaction((meas)=>{
        let res = []
        for(const m of meas){
            console.log(m.ProfTime, m.RepRate, m.ProfLen, m.ProfCnt)
            const r = stmt.run(m.ProfTime.toISOString(), m.StopDate.toISOString(), m.RepRate, m.ProfLen, m.ProfCnt, m.ProfDataDat, m.ProfDataDak, expid)
            res.push(r)
        }
        return res
    })
    
    return insertMeas(meas)
} 


module.exports = {getMeasurementsForExperimentId, insertMeasurements}
