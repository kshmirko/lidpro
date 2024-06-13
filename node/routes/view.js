var express = require('express');
var router = express.Router();
const {getExperiments} = require("../dbapi/experiment_api")

/* GET home page. */
router.get('/', async function(req, res, next) {
  const data = await getExperiments()
  
  res.render('view', { exp: data });
});

module.exports = router;
