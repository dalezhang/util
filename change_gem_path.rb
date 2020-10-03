# 先运行  bundle install --path=vendor/bundle
ENV['GEM_HOME'] = "./vendor/bundle/ruby/2.7.0"
ENV['GEM_PATH'] = "./vendor/bundle/ruby/2.7.0"
Gem.clear_paths
require "sinatra"
puts Sinatra::Application.method(:new).source_location