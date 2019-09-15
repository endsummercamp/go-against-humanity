var gulp = require('gulp'),
    plumber = require('gulp-plumber'),
    rename = require('gulp-rename');
var autoprefixer = require('gulp-autoprefixer');
var babel = require('gulp-babel');
var concat = require('gulp-concat');
var imagemin = require('gulp-imagemin'),
    cache = require('gulp-cache');
var minifycss = require('gulp-minify-css');
var less = require('gulp-less');

gulp.task('images', function(){
  return gulp.src('src/images/**/*')
    .pipe(cache(imagemin({ optimizationLevel: 3, progressive: true, interlaced: true })))
    .pipe(gulp.dest('public/images/'));
});

gulp.task('styles', function(){
  return gulp.src(['src/styles/**/*.less'])
    .pipe(plumber({
      errorHandler: function (error) {
        console.log(error.message);
        this.emit('end');
    }}))
    .pipe(less())
    .pipe(autoprefixer('last 2 versions'))
    .pipe(gulp.dest('public/styles/'))
    .pipe(rename({suffix: '.min'}))
    .pipe(minifycss())
    .pipe(gulp.dest('public/styles/'))
});

gulp.task('scripts', gulp.parallel(
  // Copy deps directly, without processing
  () => gulp.src('src/scripts/deps/*.js')
    .pipe(plumber({
      errorHandler: function (error) {
        console.log(error.message);
        this.emit('end');
    }}))
    .pipe(gulp.dest('public/scripts/')),
  // Process files in the root directory
  () => gulp.src('src/scripts/*.js')
    .pipe(plumber({
      errorHandler: function (error) {
        console.log(error.message);
        this.emit('end');
    }}))
    .pipe(babel({
      presets: ['@babel/env', '@babel/react'/*, 'minify'*/]
    }))
    .pipe(rename({suffix: '.min'}))
    .pipe(gulp.dest('public/scripts/'))
));

gulp.task('default', function(){
  gulp.watch("src/styles/**/*.less", gulp.series('styles'));
  gulp.watch("src/scripts/**/*.js", gulp.series('scripts'));
});
