class PlistsController < ApplicationController
  def new
    @plist = LaunchdPlist.new
  end

  def create
    plist = LaunchdPlist.new params[:plist]
    render :xml => LaunchdSerializer.new(plist).to_plist
  end
end
