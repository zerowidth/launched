require "spec_helper"

describe LaunchdPlist do

  let :plist do
    LaunchdPlist.new(
      :name => "Hello World",
      :command => 'growlnotify -m "hello!"',
      :weekdays => [1,2,3,4,5],
      :months => [1,4,7,10]
    )
  end

  describe "#label" do
    it "returns a computer-readable version of the name" do
      plist.label.should == "hello_world"
    end
  end

  describe "weekday_list" do
    it "returns a string listing the selected weekdays" do
      plist.weekday_list.should == "1,2,3,4,5"
    end

    it "accepts a list of weekdays as a string" do
      plist.weekday_list = "2,3,4"
      plist.weekdays.should == [2,3,4]
    end
  end

  describe "month_list" do
    it "returns a string listing the selected months" do
      plist.month_list.should == "1,4,7,10"
    end

    it "accepts a list of months as a string" do
      plist.month_list = "2,3,4"
      plist.months.should == [2,3,4]
    end
  end

end
