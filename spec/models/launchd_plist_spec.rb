require "spec_helper"

describe LaunchdPlist do

  let :plist do
    LaunchdPlist.new(
      :name => "Hello World",
      :command => 'growlnotify -m "hello!"',
      :weekdays => "1,2,3,4,5",
      :months => "1,4,7,10"
    )
  end

  describe "when saving to the database" do
    it "creates a record in the db" do
      lambda { plist.save! }.should change(LaunchdPlist, :count).by(1)
    end

    it "generates a UUID for the plist" do
      plist.save!
      plist.uuid.should_not be_nil
    end

    it "disallows invalid minute values" do
      plist.minute = "iamaspammer"
      plist.save
      plist.should have(1).error_on(:minute)
    end
  end

  describe "#label" do
    it "returns a computer-readable version of the name" do
      plist.label.should == "hello_world"
    end
  end

  describe "#month_list" do
    it "returns an integer array of months" do
      plist.month_list.should == [1,4,7,10]
    end

    it "returns an empty array when months is nil" do
      plist.months = nil
      plist.month_list.should == []
    end
  end

  describe "weekday_list" do
    it "returns an integer array of weekdays" do
      plist.weekday_list.should == [1,2,3,4,5]
    end

    it "returns an empty array when weekdays is nil" do
      plist.weekdays = nil
      plist.weekday_list.should == []
    end
  end

end
