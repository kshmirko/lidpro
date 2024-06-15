var express = require('express');
const { getMeasurementsForExperimentId } = require('../dbapi/measurement_api');
const { getExperimentById } = require('../dbapi/experiment_api');
var router = express.Router();

/* GET home page. */
router.get('/', async function (req, res, next) {
    res.render('process');
});



router.get('/experiment/:expid', async function (req, res, next) {
    const exp_id = req.params.expid
    const meas = getMeasurementsForExperimentId(exp_id)

    let total = {
        StartTime: new Date(),
        StopTime: new Date(2010, 1, 1, 12, 0, 0),
        ProfCnt: 0
    }
    let dat = new Uint32Array(meas[0].ProfLen)
    let dak = new Uint32Array(meas[0].ProfLen)

    for (let m of meas) {
        if (total.StartTime > m.ProfTime) {
            total.StartTime = m.ProfTime
        }
        if (total.StopTime < m.ProfTime) {
            total.StopTime = m.ProfTime
        }

        for (let j = 0; j < dat.length; j++) {
            dat[j] += m.Dat[j]
            dak[j] += m.Dak[j]
        }
        total.ProfCnt += m.ProfCnt
    }

    total.Dat = dat;
    total.Dak = dak;

    
    res.render('process');
});

router.get('/:expid/:mid', async function (req, res, next) {
    const exp_id = req.params.expid
    const m_id = req.params.mid
    res.render('process');
});

module.exports = router;
