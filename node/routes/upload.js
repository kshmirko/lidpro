var express = require('express');
const fs = require("fs");
const multer  = require('multer')

//const storage = multer.memoryStorage()
//const upload = multer({ storage: storage })
const upload = multer({dest:"./uploads"})
var router = express.Router();
const {db} = require("../db")
const {insertExperiment} = require("../dbapi/experiment_api")
const {insertMeasurements} = require("../dbapi/measurement_api")
const {readLidarZipBuffer} = require("../utils/lidarfile")



/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('upload', { title: 'Upload' });
});


router.post('/', upload.single('experiment-archivefile'), async function(req, res, next) {
    const data = Buffer.from(
                      fs.readFileSync(
                        req.file.path, {
                          encoding:'binary',
                          flag:'r'
                        }),
                        'binary')
    
    const rec = {
        StartTime: req.body['experiment-datetime'].replace("T", " "),
        Title: req.body['experiment-title'],
        Comment: req.body['experiment-comment'],
        VertRes: req.body['experiment-vertres'],
        AccumTime: req.body['experiment-accumtime'],
        Archive: data,
    }

    const p = await readLidarZipBuffer(req.file.path, rec.AccumTime)
    const insertedId1 = insertExperiment(rec)
    const insertedId2 = insertMeasurements(p, insertedId1)


    

    res.render('upload', {status:'Данные эксперимента успешно згружены  '})
})

module.exports = router;
