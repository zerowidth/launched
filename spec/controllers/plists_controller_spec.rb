require 'spec_helper'

describe PlistsController do
  describe "POST to create" do

    def create(overrides={})
      post :create,
        :format => :xml,
        :plist => {
          :name => "test command",
          :command => "echo hello",
          :interval => "300"
        }.merge(overrides)
    end

    it "redirects to the new plist path, using uuid as the id" do
      create
      plist = LaunchdPlist.last
      response.should redirect_to(plist_path(plist.uuid, :format => :xml))
    end

    it "creates a launchd plist entry" do
      lambda { create }.should change(LaunchdPlist, :count).by(1)
    end
  end

  describe "GET to show with a UUID" do

    let :plist do
      LaunchdPlist.create! :name => "test", :command => "ls", :interval => "300"
    end

    before do
      get :show, :id => plist.uuid
    end

    it "renders an xml plist" do
      response.body.should =~ /xml.*StartInterval/m
    end
  end
end
