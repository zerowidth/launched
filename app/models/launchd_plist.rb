class LaunchdPlist < ActiveRecord::Base

  CRON_EXP = /^[\-\/*0-9,]+$/

  validates_presence_of :uuid, :command, :name
  validates_format_of :minute,
    :with => CRON_EXP, :allow_nil => true, :allow_blank => true
  validates_format_of :hour,
    :with => CRON_EXP, :allow_nil => true, :allow_blank => true
  validates_format_of :day_of_month,
    :with => CRON_EXP, :allow_nil => true, :allow_blank => true

  before_validation :generate_uuid, :on => :create

  def label
    name.downcase.gsub /\s+/, "_"
  end

  def weekday_list
    (weekdays || "").split(",").map(&:to_i)
  end

  def month_list
    (months || "").split(",").map(&:to_i)
  end

  protected

  def generate_uuid
    self.uuid = UUID.generate
  end

end
