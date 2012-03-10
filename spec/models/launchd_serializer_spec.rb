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
        plist.weekdays = "1,5"
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

    context "when the plist is safe to launch only once" do
      it "sets LaunchOnlyOnce to true" do
        plist.launch_only_once = true
        xml["LaunchOnlyOnce"].should == true
      end
    end

    context "when a user is set" do
      it "sets UserName" do
        plist.user = "bobby"
        xml["UserName"].should == "bobby"
      end

      it "does not set UserName when user is an empty string" do
        plist.user = ""
        xml["UserName"].should == nil
      end
    end

    context "when a group is set" do
      it "sets GroupName" do
        plist.group = "wheel"
        xml["GroupName"].should == "wheel"
      end

      it "does not set GroupName with an empty string" do
        plist.group = ""
        xml["GroupName"].should == nil
      end
    end

    context "with a root directory set" do
      it "sets RootDirectory" do
        plist.root_directory = "/tmp"
        xml["RootDirectory"].should == "/tmp"
      end

      it "deos not set RootDirectory when root directory is a blank string" do
        plist.root_directory = ""
        xml["RootDirectory"].should == nil
      end
    end

    context "with a working directory set" do
      it "sets WorkingDirectory" do
        plist.working_directory = "/tmp"
        xml["WorkingDirectory"].should == "/tmp"
      end

      it "does not set WorkingDirectory when working directory is a blank string" do
        plist.working_directory = ""
        xml["WorkingDirectory"].should == nil
      end
    end
  end

end
