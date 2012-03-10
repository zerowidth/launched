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
    @plist = LaunchdPlist.find_by_uuid(params[:id]) or raise ActiveRecord::RecordNotFound
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
end
