let mix = require('laravel-mix');
require('laravel-mix-purgecss');

mix.setPublicPath('./public');
mix.sass('resources/assets/sass/app.scss', 'public/css')
   .version();

