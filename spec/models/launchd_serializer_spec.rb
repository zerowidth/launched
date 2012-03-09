require "spec_helper"

describe LaunchdSerializer do

  let :plist do
    LaunchdPlist.new(
      :name => "Hello World",
      :command => 'growlnotify -m "hello!"'
    )
  end

  let :serializer do
    LaunchdSerializer.new(plist)
  end

  let :xml do
    Plist.parse_xml serializer.to_plist
  end

  describe "#to_plist" do
    it "returns an xml string" do
      serializer.to_plist.should =~ /xml/
    end

    it "includes the job name as a fully qualified label" do
      xml["Label"].should == "com.zerowidth.launched.hello_world"
    end

    it "includes the command in ProgramArguments" do
      args = xml["ProgramArguments"]
      args.should have(3).items
      args.last.should =~ /growlnotify/
    end

    context "with a single minute specified" do
      it "sets a single StartCalendarInterval" do
        plist.minute = "10"
        xml["StartCalendarInterval"].should == {
          "Minute" => 10
        }
      end
    end

    context "with a single minute and multiple hours specified" do
      it "sets multiple StartCalendarIntervals" do
        plist.minute = "0"
        plist.hour = "8,17"
        xml["StartCalendarInterval"].should == [
          {"Minute" => 0, "Hour" => 8},
          {"Minute" => 0, "Hour" => 17}
        ]
      end
    end

    context "with a weekday list specified" do
      it "sets the weekdays to run" do
        plist.weekday_list = "1,5"
        xml["StartCalendarInterval"].should == [
          {"Weekday" => 1},
          {"Weekday" => 5}
        ]
      end
    end

    context "with a run interval" do
      it "sets the StartInterval key" do
        plist.interval = 300
        xml["StartInterval"].should == 300
      end
    end

    context "with run at load set" do
      it "sets the RunAtLoad boolean" do
        plist.run_at_load = true
        xml["RunAtLoad"].should == true
      end
    end
  end

end
