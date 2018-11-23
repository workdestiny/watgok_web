let mix = require('laravel-mix');
require('laravel-mix-purgecss');

mix.setPublicPath('./public');
mix.js('resources/assets/js/app.js', 'public/js')
   .sass('resources/assets/sass/app.scss', 'public/css')
   .version();