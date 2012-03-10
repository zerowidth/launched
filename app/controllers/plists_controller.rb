class PlistsController < ApplicationController
  def new
    @plist = LaunchdPlist.new
  end

  def create
    @plist = LaunchdPlist.new params[:plist]

    if @plist.save
      redirect_to plist_path(@plist.uuid, :format => :xml)
    else
      render :action => :new
    end
  end

  def show
    plist = LaunchdPlist.find_by_uuid params[:id]
    render :xml => LaunchdSerializer.new(plist).to_plist
  end
end
