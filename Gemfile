source 'https://rubygems.org'

ruby "2.7.2"

gem 'rake'
gem 'rails'
gem "webpacker"

# Reduces boot times through caching; required in config/boot.rb
gem 'bootsnap', '>= 1.4.4', require: false

gem "pg"
gem "uuid"
gem "plist"
# gem 'anjlab-bootstrap-rails', :require => 'bootstrap-rails',
#                               :git => 'git://github.com/anjlab/bootstrap-rails.git'

group :development, :test do
  # Call 'byebug' anywhere in the code to stop execution and get a debugger console
  gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]
end

group :development do
  # Access an interactive console on exception pages or by calling 'console' anywhere in the code.
  gem 'web-console', '>= 4.1.0'
  # Display performance information such as SQL time and flame graphs for each request in your browser.
  # Can be configured to work on production as well see: https://github.com/MiniProfiler/rack-mini-profiler/blob/master/README.md
  # gem 'rack-mini-profiler', '~> 2.0'
  gem 'listen', '~> 3.3'
  # Spring speeds up development by keeping your application running in the background. Read more: https://github.com/rails/spring
  gem 'spring'
end

group :test do
  gem "rails-controller-testing"
  gem "rspec-rails"
end

# Windows does not include zoneinfo files, so bundle the tzinfo-data gem
gem 'tzinfo-data', platforms: [:mingw, :mswin, :x64_mingw, :jruby]

# gem 'sass-rails', "~> 5.0.7"
# gem 'coffee-rails', '~> 4.2.2'
# gem 'uglifier', '>= 1.0.3'
# gem 'jquery-rails'
