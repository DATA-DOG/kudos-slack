var gulp = require('gulp');
var exec = require('child_process').exec;
var watch = require('gulp-watch');

gulp.task('build', function (cb) {
  exec('go build', function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
})

gulp.task('run', function(cb) {
  exec('./kudos', function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
});

gulp.task('chmod', function(cb) {
  exec('chmod +x kudos', function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
});

gulp.task('default', ['build', 'chmod'], function() {
  gulp.watch('*.go', ['build']);
})
