class PlistsController < ApplicationController

  before_filter :find_plist_by_uuid, :only => [:show, :edit, :install]

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
    @plist_xml = LaunchdSerializer.new(@plist).to_plist

    respond_to do |format|
      format.xml do
        filename = Launched::Application::DOMAIN + "." + @plist.label + ".xml"
        headers["Content-Disposition"] =
          %Q(attachment; filename="#{filename}"; size=#{@plist_xml.length})
        render :xml => @plist_xml
      end
      format.html
    end
  end

  def edit
    render :action => "new"
  end

  def install
    render :template => "plists/install.txt"
  end

  protected

  def find_plist_by_uuid
    @plist = LaunchdPlist.find_by_uuid(params[:id]) or raise ActiveRecord::RecordNotFound
  end
end
