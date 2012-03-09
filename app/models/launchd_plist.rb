class LaunchdPlist

  attr_accessor :command, :name
  attr_accessor :minute, :hour, :day_of_month, :weekdays, :months
  attr_reader :interval, :run_at_load, :launch_only_once

  def initialize(attributes={})
    attributes.each do |k, v|
      self.send "#{k}=", v
    end
    self.weekdays ||= []
    self.months ||= []
  end

  def interval=(value)
    @interval = value && !value.blank? ? value.to_i : nil
  end

  def run_at_load=(value)
    @run_at_load = value == "1" || value == true
  end

  def launch_only_once=(value)
    @launch_only_once = value == "1" || value == true
  end

  def label
    name.downcase.gsub /\s+/, "_"
  end

  def weekday_list
    weekdays.join ","
  end

  def weekday_list=(list)
    @weekdays = list.split(",").map(&:to_i)
  end

  def month_list
    months.join ","
  end

  def month_list=(list)
    @months = list.split(",").map(&:to_i)
  end

end
