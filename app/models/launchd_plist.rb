class LaunchdPlist
  include ActiveModel::Model
  include ActiveModel::Attributes
  include ActiveModel::Serializers::JSON

  attribute :uuid, :string

  # basics
  attribute :name, :string
  attribute :command, :string

  # scheduling
  attribute :start_interval, :integer
  attribute :minute, :string
  attribute :hour, :string
  attribute :day_of_month, :string
  attribute :month, :string
  attribute :weekday, :string

  # runtime
  attribute :user, :string
  attribute :group, :string
  attribute :root_directory, :string
  attribute :working_directory, :string

  # just for tracking
  attribute :created_at, :datetime

  CRON_EXP = %r{\A[\-/*0-9,]+\z}

  validates :command, presence: true, length: { maximum: 1024 }
  validates :name, presence: true, length: { maximum: 256 }
  validates :start_interval, numericality: true, allow_blank: true
  validates :minute, :hour, :day_of_month, :month, :weekday,
    length: { maximum: 64 }, format: { with: CRON_EXP }, allow_blank: true
  validates :user, :group, :root_directory, :working_directory,
    length: { maximum: 256 }, allow_blank: true

  def self.count
    REDIS.with do |redis|
      redis.scan_each(match: "#{namespace}:*").count
    end
  end

  def self.all
    REDIS.with do |redis|
      redis.scan_each(match: "#{namespace}:*").map do |key|
        find key.split(":")[1]
      end
    end
  end

  def self.namespace
    name.underscore
  end

  def self.find(uuid)
    REDIS.with do |redis|
      if (json = redis.get("#{namespace}:#{uuid}"))
        plist = new
        plist.from_json(json)
        plist.uuid = uuid
        plist.persisted = true
        plist
      end
    end
  end

  def initialize(...)
    super
    self.created_at ||= Time.now.utc
    self.uuid ||= Digest::UUID.uuid_v4
    self.persisted = false
  end

  attr_accessor :persisted

  def persisted?
    persisted
  end

  def new_record?
    !persisted
  end

  def save
    return false unless valid?

    REDIS.with do |redis|
      redis.set("#{self.class.namespace}:#{uuid}", to_json)
    end
    self.persisted = true
    true
  end

  # For ActiveModel::Serialization
  def attributes
    {
      "command" => command,
      "name" => name,
      "start_interval" => start_interval,
      "minute" => minute,
      "hour" => hour,
      "day_of_month" => day_of_month,
      "month" => month,
      "weekday" => weekday,
      "user" => user,
      "group" => group,
      "root_directory" => root_directory,
      "working_directory" => working_directory,
      "created_at" => created_at
    }
  end

  # for from_json support
  def attributes=(hash)
    hash.each do |key, value|
      send("#{key}=", value)
    end
  end

  def label
    name&.downcase&.gsub(/\s+/, "_")
  end
end
