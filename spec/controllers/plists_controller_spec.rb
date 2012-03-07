require 'spec_helper'

describe PlistsController do
  describe "POST to create" do
    before do
      post :create,
        :format => :xml,
        :plist => {
          :name => "test command",
          :command => "echo hello",
          :interval => "300"
        }
    end

    it "is successful" do
      response.should be_success
    end

    it "renders an xml plist" do
      response.body.should =~ /xml.*StartInterval/m
    end
  end
end
