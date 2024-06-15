var express = require('express');
var router = express.Router();
const {getExperiments, getExperimentById} = require("../dbapi/experiment_api")
const {getMeasurementsForExperimentId} = require("../dbapi/measurement_api")

/* GET home page. */
router.get('/', async function(req, res, next) {
  const data = await getExperiments()
  
  res.render('view', { exp: data });
});

/* GET home page. */
router.get('/experiment/:id', async function(req, res, next) {

  const id = req.params.id
  const exp = getExperimentById(id)
  
  const meas = getMeasurementsForExperimentId(id)
  
  res.render('detail', { 
    exp:exp, 
    meas: meas 
  });
});



module.exports = router;
