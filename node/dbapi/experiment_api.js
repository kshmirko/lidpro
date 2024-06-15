const {db} = require("./../db")

/**
 * 
 * @returns {List<Object>} - список объектов "Эксперимент"
 */
function getExperiments() {
    let experiments = []
    
    const stmt = db.prepare(`
    SELECT 
        * 
    FROM Experiment
    ORDER BY
        START_TIME ASC;
    `)
    for(const exp of stmt.iterate()){
        experiments.push({
            Id: exp.ID,
            StartTime: exp.START_TIME,
            Title: exp.TITLE,
            Comment: exp.COMMENT,
            VertRes: exp.VERT_RES,
            AccumTime: exp.ACCUM_TIME,
            Archive:exp.ARCHIVE || ''
        })
    }
    return  experiments
}

/**
 * 
 * @param {number} id - идентификатор эксперимента, целое число 
 * @returns {Object} - структура, содержащая параметры эксперимента
 */
function getExperimentById(id) {
    const stmt = db.prepare(`SELECT * FROM Experiment WHERE ID = ?`)
    exp = stmt.get(id)

    return {
        Id: exp.ID,
        StartTime: exp.START_TIME,
        Title: exp.TITLE,
        Comment: exp.COMMENT,
        VertRes: exp.VERT_RES,
        AccumTime: exp.ACCUM_TIME,
        Archive: exp.ARCHIVE || ''
    }
}


/**
 * 
 * @param {Object} rec - структура, содержащая параметры эксперимента 
 * @returns {number} - идентификатор добавленной записи
 */
function insertExperiment(rec) {
    const stmt = db.prepare(`
    INSERT INTO 
        Experiment(START_TIME, TITLE, COMMENT, VERT_RES, ACCUM_TIME, ARCHIVE)
    VALUES(?,?,?,?,?,?);
    `)

    //let res
    const insertExp = db.transaction((r)=>{
        const res = stmt.run(r.StartTime, r.Title, r.Comment, r.VertRes, r.AccumTime, r.Archive)
        return res
    })
    const ret = insertExp(rec)
    return ret.lastInsertRowid
}


module.exports = {getExperiments, getExperimentById, insertExperiment}