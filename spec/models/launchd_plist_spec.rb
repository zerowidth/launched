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

  describe "interval" do
    it "sets the interval to an integer when assigned a string" do
      plist.interval = "300"
      plist.interval.should == 300
    end

    it "sets the interval to nil when assigned an empty string" do
      plist.interval = ""
      plist.interval.should == nil
    end
  end

  describe "run at load" do
    it "is false when a blank string is assigned" do
      plist.run_at_load = ""
      plist.run_at_load.should == false
    end

    it "is false when a '0' is assigned" do
      plist.run_at_load = "0"
      plist.run_at_load.should == false
    end

    it "is true when a '1' is assigned" do
      plist.run_at_load = "1"
      plist.run_at_load.should == true
    end
  end

  describe "launch only once" do
    it "is false when a blank string is assigned" do
      plist.launch_only_once = ""
      plist.launch_only_once.should == false
    end

    it "is false when a '0' is assigned" do
      plist.launch_only_once = "0"
      plist.launch_only_once.should == false
    end

    it "is true when a '1' is assigned" do
      plist.launch_only_once = "1"
      plist.launch_only_once.should == true
    end
  end

end
