require('dotenv').config()

var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');


const {db} = require("./db")

// Load all necessary routes from ajacent files 
var indexRouter = require('./routes/index');
var uploadRouter= require('./routes/upload')
var viewRouter = require('./routes/view')
var operateRouter = require("./routes/operate")

var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

// setup logger
app.use(logger('dev'));

app.use(express.json());
app.use(express.urlencoded({ extended: false }));
//app.use(express.multipart())
app.use(cookieParser());
// setup static dir
app.use(express.static(path.join(__dirname, 'public')));

// mount routes
app.use('/', indexRouter);
app.use('/api/upload',uploadRouter)
app.use('/api/view', viewRouter)
app.use("/api/process", operateRouter)

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

module.exports = app;
