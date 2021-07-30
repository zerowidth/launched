require "rails_helper"

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
      expect(serializer.to_plist).to match(/xml/)
    end

    it "includes the job name as a fully qualified label" do
      expect(xml["Label"]).to eq("com.zerowidth.launched.hello_world")
    end

    it "includes the command in ProgramArguments" do
      args = xml["ProgramArguments"]
      expect(args.length).to be(3)
      expect(args.last).to match(/growlnotify/)
    end

    context "with a single minute specified" do
      it "sets a single StartCalendarInterval" do
        plist.minute = "10"
        expect(xml["StartCalendarInterval"]).to eq("Minute" => 10)
      end
    end

    context "with a single minute and multiple hours specified" do
      it "sets multiple StartCalendarIntervals" do
        plist.minute = "0"
        plist.hour = "8,17"
        expect(xml["StartCalendarInterval"]).to match_array [
          {"Minute" => 0, "Hour" => 8},
          {"Minute" => 0, "Hour" => 17}
        ]
      end
    end

    context "with a weekday list specified" do
      it "sets the weekdays to run" do
        plist.weekdays = "1,5"
        expect(xml["StartCalendarInterval"]).to match_array([
          {"Weekday" => 1},
          {"Weekday" => 5}
        ])
      end
    end

    context "with a run interval" do
      it "sets the StartInterval key" do
        plist.interval = 300
        expect(xml["StartInterval"]).to be(300)
      end
    end

    context "with run at load set" do
      it "sets the RunAtLoad boolean" do
        plist.run_at_load = true
        expect(xml["RunAtLoad"]).to be(true)
      end
    end

    context "when the plist is safe to launch only once" do
      it "sets LaunchOnlyOnce to true" do
        plist.launch_only_once = true
        expect(xml["LaunchOnlyOnce"]).to be(true)
      end
    end

    context "when a user is set" do
      it "sets UserName" do
        plist.user = "bobby"
        expect(xml["UserName"]).to eq("bobby")
      end

      it "does not set UserName when user is an empty string" do
        plist.user = ""
        expect(xml["UserName"]).to be_nil
      end
    end

    context "when a group is set" do
      it "sets GroupName" do
        plist.group = "wheel"
        expect(xml["GroupName"]).to eq "wheel"
      end

      it "does not set GroupName with an empty string" do
        plist.group = ""
        expect(xml["GroupName"]).to be_nil
      end
    end

    context "with a root directory set" do
      it "sets RootDirectory" do
        plist.root_directory = "/tmp"
        expect(xml["RootDirectory"]).to eq "/tmp"
      end

      it "deos not set RootDirectory when root directory is a blank string" do
        plist.root_directory = ""
        expect(xml["RootDirectory"]).to be_nil
      end
    end

    context "with a working directory set" do
      it "sets WorkingDirectory" do
        plist.working_directory = "/tmp"
        expect(xml["WorkingDirectory"]).to eq "/tmp"
      end

      it "does not set WorkingDirectory when working directory is a blank string" do
        plist.working_directory = ""
        expect(xml["WorkingDirectory"]).to be_nil
      end
    end
  end

end
