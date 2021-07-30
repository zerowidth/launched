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
      expect { plist.save! }.to change(LaunchdPlist, :count).by(1)
    end

    it "generates a UUID for the plist" do
      plist.save!
      expect(plist.uuid).to_not be_nil
    end

    it "disallows invalid minute values" do
      plist.minute = "iamaspammer"
      expect(plist).to_not be_valid
      expect(plist.errors[:minute]).to be_present
    end
  end

  describe "#label" do
    it "returns a computer-readable version of the name" do
      expect(plist.label).to eq("hello_world")
    end
  end

  describe "#month_list" do
    it "returns an integer array of months" do
      expect(plist.month_list).to eq [1,4,7,10]
    end

    it "returns an empty array when months is nil" do
      plist.months = nil
      expect(plist.month_list).to be_empty
    end
  end

  describe "weekday_list" do
    it "returns an integer array of weekdays" do
      expect(plist.weekday_list).to eq [1,2,3,4,5]
    end

    it "returns an empty array when weekdays is nil" do
      plist.weekdays = nil
      expect(plist.weekday_list).to be_empty
    end
  end

end
