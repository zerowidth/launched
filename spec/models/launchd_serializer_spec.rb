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
      expect(xml["Label"]).to eq("#{Launched::DOMAIN}.hello_world")
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
        plist.weekday = "1,5"
        expect(xml["StartCalendarInterval"]).to match_array([
          {"Weekday" => 1},
          {"Weekday" => 5}
        ])
      end
    end

    context "with a start interval" do
      it "sets the StartInterval key" do
        plist.start_interval = 300
        expect(xml["StartInterval"]).to be(300)
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

    context "with logging paths" do
      it "includes StandardOutPath and StandardErrorPath" do
        plist.standard_out_path = "/tmp/stdout.log"
        plist.standard_error_path = "/tmp/stderr.log"
        expect(xml["StandardOutPath"]).to eq "/tmp/stdout.log"
        expect(xml["StandardErrorPath"]).to eq "/tmp/stderr.log"
      end
    end
  end

end
