require 'spec_helper'

describe PlistsController do

  let :plist do
    LaunchdPlist.create! :name => "test", :command => "ls", :interval => "300"
  end

  describe "POST to create" do

    def create(overrides={})
      post :create,
        :plist => {
          :name => "test command",
          :command => "echo hello",
          :interval => "300"
        }.merge(overrides)
    end

    it "redirects to the new plist path, using uuid as the id" do
      create
      plist = LaunchdPlist.last
      response.should redirect_to(plist_path(plist.uuid))
    end

    it "creates a launchd plist entry" do
      lambda { create }.should change(LaunchdPlist, :count).by(1)
    end
  end

  describe "GET to show with a UUID" do

    context "when requesting xml" do
      before do
        get :show, :id => plist.uuid, :format => :xml
      end

      it "renders an xml plist with the xml format" do
        response.body.should =~ /xml.*StartInterval/m
      end

      it "sets the content-disposition header" do
        response.headers["Content-Disposition"].should =~ /attachment.*test\.xml/
      end
    end

    it "renders html with the html format" do
      get :show, :id => plist.uuid
      response.should be_success
      response.should render_template("show")
    end

    it "raises a record not found error for invalid ids" do
      lambda do
        get :show, :id => 'lol'
      end.should raise_error(ActiveRecord::RecordNotFound)
    end

  end

  describe "GET to edit with a UUID" do
    it "is succesful" do
      get :edit, :id => plist.uuid
      response.should be_success
    end

    it "renders the new plist form" do
      get :edit, :id => plist.uuid
      response.should render_template("new")
    end

    it "raises a record not found error for invalid ids" do
      lambda do
        get :edit, :id => 'lol'
      end.should raise_error(ActiveRecord::RecordNotFound)
    end
  end

  describe "GET to install with a UUID" do
    it "is successful" do
      get :install, :id => plist.uuid
      response.should be_success
    end
  end
end
