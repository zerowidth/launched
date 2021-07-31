require "rails_helper"

describe PlistsController do
  before do
    REDIS.with do |redis|
      redis.scan_each(match: "#{LaunchdPlist.namespace}:*").each do |key|
        redis.del key
      end
    end
  end

  let :plist do
    p = LaunchdPlist.new(name: "test", command: "ls", start_interval: "300")
    unless p.save
      raise p.errors.inspect
    end
    p
  end

  describe "POST to create" do
    def create(overrides = {})
      post :create,
        params: { plist: {
          name: "test command",
          command: "echo hello",
          start_interval: "300",
        }.merge(overrides) }
    end

    it "redirects to the new plist path" do
      create
      plist = LaunchdPlist.all.first
      expect(response).to redirect_to(plist_path(plist.uuid))
    end

    it "creates a launchd plist entry" do
      expect(lambda { create }).to change(LaunchdPlist, :count).by(1)
    end
  end

  describe "GET to show with a UUID" do
    context "when requesting xml" do
      before do
        get :show, params: { id: plist.uuid }, format: :xml
      end

      it "renders an xml plist with the xml format" do
        expect(response.body).to match(/xml.*StartInterval/m)
      end

      it "sets the content-disposition header" do
        expect(response.headers["Content-Disposition"]).to match(/attachment.*test\.xml/)
      end
    end

    it "renders html with the html format" do
      get :show, params: { id: plist.uuid }
      expect(response).to be_successful
      expect(response).to render_template("show")
    end

    it "raises a record not found error for invalid ids" do
      expect { get :show, params: { id: "lol" }, format: :xml }.to raise_error(ActionController::RoutingError)
    end
  end

  describe "GET to edit with a UUID" do
    it "is succesful" do
      get :edit, params: { id: plist.uuid }
      expect(response).to be_successful
    end

    it "renders the new plist form" do
      get :edit, params: { id: plist.uuid }
      expect(response).to render_template("new")
    end

    it "raises an error for invalid ids" do
      expect { get :edit, params: { id: "lol" } }.to raise_error(ActionController::RoutingError)
    end
  end

  describe "GET to install with a UUID" do
    it "is successful" do
      get :install, params: { id: plist.uuid }
      expect(response).to be_successful
    end
  end
end
