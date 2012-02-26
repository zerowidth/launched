require 'spec_helper'

describe CommandPlist do
  context "when initialized with a simple command with no schedule" do

    let(:command) { Command.new :name => "Hello World", :command => "ls -al" }
    let(:plist)   { CommandPlist.new command }
    let(:parsed)  { Plist.parse_xml(plist.plist) }

    describe "#plist" do
      it "generates an xml string" do
        plist.plist.should =~ /<?xml/
      end

      it "sets the command name as the label" do
        parsed["Label"].should == "com.zerowidth.launched.hello_world"
      end

      it "sets the command as part of the ProgramArguments list" do
        parsed["ProgramArguments"].should == [
          "sh",
          "-c",
          "ls -al"
        ]
      end
    end

  end
end
