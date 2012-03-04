class PlistsController < ApplicationController
  def new
    @plist = LaunchdPlist.new
  end

  def create
  end
end
