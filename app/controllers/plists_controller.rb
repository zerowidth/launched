class PlistsController < ApplicationController
  before_action :find_plist_by_uuid, only: %i[show edit install]

  def new
    @plist = LaunchdPlist.new
  end

  def create
    @plist = LaunchdPlist.new(plist_params)

    if @plist.save
      redirect_to plist_path(@plist.uuid)
    else
      render action: :new
    end
  end

  def show
    @plist_xml = LaunchdSerializer.new(@plist).to_plist

    respond_to do |format|
      format.xml do
        filename = "#{Launched::Application::DOMAIN}.#{@plist.label}.xml"
        headers["Content-Disposition"] =
          %(attachment; filename="#{filename}"; size=#{@plist_xml.length})
        render xml: @plist_xml
      end
      format.html
    end
  end

  def edit
    render action: "new"
  end

  def install
    render "plists/install", formats: [:txt]
  end

  protected

  def plist_params
    params.require(:plist).permit(:name, :command, :interval)
  end

  def find_plist_by_uuid
    @plist = LaunchdPlist.find(params[:id]) || raise(ActionController::RoutingError, "not found")
  end
end
