source 'https://rubygems.org'
git_source(:github) { |repo| "https://github.com/#{repo}.git" }

ruby '3.1.1'

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails', branch: 'main'
gem 'rails', '~> 7'
# Use Puma as the app server
gem 'puma', '~> 5.6'
# Use SCSS for stylesheets
gem 'sass-rails', '>= 6'
# Use Active Model has_secure_password
# gem 'bcrypt', '~> 3.1.7'

gem "bootstrap"
gem "connection_pool"
gem "plist"
gem "redis"
gem "simple_form"
gem "mini_portile2", "~> 2.8.0" # for nokogiri to build

group :development, :test do
  gem "debug", ">= 1.0.0"
  gem "rspec-rails"
end

group :development do
  gem 'listen', '~> 3.3'
end

group :test do
  gem "rails-controller-testing"
end

# Windows does not include zoneinfo files, so bundle the tzinfo-data gem
gem 'tzinfo-data', platforms: [:mingw, :mswin, :x64_mingw, :jruby]
