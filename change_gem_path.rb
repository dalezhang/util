# 先运行  bundle install --path=vendor/bundle
# In Windows it is set
#
# set GEM_HOME=[path]/projects/shared/gems/ruby/1.8/gems
# Linux would be export
#
# export GEM_HOME=~/projects/shared/gems/ruby/1.8/gems

ENV['GEM_HOME'] = "./vendor/bundle/ruby/2.7.0"
ENV['GEM_PATH'] = "./vendor/bundle/ruby/2.7.0"
Gem.clear_paths
require "byebug"
puts self.method(:byebug).source_location