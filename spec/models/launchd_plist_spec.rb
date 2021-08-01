require "spec_helper"

describe LaunchdPlist do
  before do
    REDIS.with do |redis|
      redis.scan_each(match: "#{LaunchdPlist.namespace}:*").each do |key|
        redis.del key
      end
    end
  end

  let :plist do
    LaunchdPlist.new(
      :name => "Hello World",
      :command => 'growlnotify -m "hello!"',
      :weekday => "1,2,3,4,5",
      :month => "1,4,7,10"
    )
  end

  let(:active_model_instance) { plist }
  it_behaves_like "ActiveModel"

  it "is valid" do
    expect(plist).to be_valid
  end

  describe "attributes" do
    it "leaves off nil entries" do
      expect(plist.attributes.keys).to match_array(%w[name command weekday month created_at])
    end

    it "leaves off empty entries" do
      expect(LaunchdPlist.new(name: "").attributes.keys).to_not include("name")
    end
  end

  describe "when saving to the database" do
    it "creates a record in the db" do
      expect { plist.save }.to change(LaunchdPlist, :count).by(1)
    end

    it "retrieves a previously saved record" do
      attributes = {
        name: "testing",
        command: "echo hello world",
        start_interval: 10,
        minute: "0",
        hour: "12",
        day_of_month: "*/2",
        month: "*",
        weekday: "1,3,5",
        user: "whoami",
        group: "wheel",
        root_directory: "/Users/whoami",
        working_directory: "/tmp",
        standard_out_path: "/var/log/stdout.log",
        standard_error_path: "/var/log/stderr.log",
      }.stringify_keys
      plist = LaunchdPlist.new(attributes)
      expect(plist.save).to be true
      from_db = LaunchdPlist.find(plist.uuid)
      expect(from_db).to be_present
      expect(from_db.attributes.except("created_at")).to eq(attributes.except("created_at"))
      expect(from_db.created_at.to_i).to eq(plist.created_at.to_i)
    end

    it "generates a UUID for the plist" do
      expect(plist.save).to be true
      expect(plist.uuid&.length).to eq(36)
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
end
