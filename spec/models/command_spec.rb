require 'spec_helper'

describe Command do
  subject { Command.new :name => "a command", :command => "ls -la" }

  it "can be initialized with nil" do
    lambda { Command.new nil }.should_not raise_error
  end

  it "is valid" do
    subject.should be_valid
  end

  it "requires a name" do
    subject.name = nil
    subject.should_not be_valid
    subject.should have(1).error_on(:name)
  end

  it "requires a command" do
    subject.command = nil
    subject.should_not be_valid
    subject.should have(1).error_on(:command)
  end
end
