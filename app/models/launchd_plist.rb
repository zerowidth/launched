class LaunchdPlist
  include ActiveModel::Model
  include ActiveModel::Attributes
  include ActiveModel::Serializers::JSON
  include ActiveModel::Validations

  attribute :uuid, :string
  attribute :command, :string # limit 1024
  attribute :name, :string
  attribute :minute, :string
  attribute :hour, :string
  attribute :day_of_month, :string
  attribute :weekdays, :string
  attribute :months, :string
  attribute :interval, :integer
  attribute :run_at_load, :boolean
  attribute :launch_only_once, :boolean
  attribute :user, :string
  attribute :group, :string
  attribute :root_directory, :string
  attribute :working_directory, :string
  attribute :created_at, :datetime

  CRON_EXP = %r{\A[\-/*0-9,]+\z}

  validates :command, presence: true, length: { maximum: 1024 }
  validates :name, presence: true, length: { maximum: 256 }
  validates :minute, :hour, :day_of_month, format: { with: CRON_EXP }, allow_blank: true
  # TODO more string limits

  def self.count
    REDIS.with do |redis|
      redis.scan_each(match: "#{namespace}:*").count
    end
  end

  def self.all
    REDIS.with do |redis|
      redis.scan_each(match: "#{namespace}:*").map do |uuid|
        find uuid.split(":")[1]
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
        plist
      end
    end
  end

  def initialize(...)
    super
    self.created_at ||= Time.now.utc
    self.uuid ||= UUID.generate
  end

  def save
    if valid?
      REDIS.with do |redis|
        redis.set("#{self.class.namespace}:#{uuid}", to_json)
      end
      true
    end
  end

  # For ActiveModel::Serialization
  def attributes
    {
      "uuid" => uuid, # redundant but it's fine
      "command" => command,
      "name" => name,
      "minute" => minute,
      "hour" => hour,
      "day_of_month" => day_of_month,
      "weekdays" => weekdays,
      "months" => months,
      "interval" => interval,
      "run_at_load" => run_at_load,
      "launch_only_once" => launch_only_once,
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
    name.downcase.gsub(/\s+/, "_")
  end

  def weekday_list
    (weekdays || "").split(",").map(&:to_i)
  end

  def month_list
    (months || "").split(",").map(&:to_i)
  end
end
