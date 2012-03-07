class LaunchdSerializer
  attr_reader :plist

  def initialize(plist)
    @plist = plist
  end

  def to_plist
    attribute_hash.to_plist
  end

  protected

  def attribute_hash
    {}.tap do |attrs|
      attrs["Label"] = Launched::Application::DOMAIN + "." + plist.label
      attrs["ProgramArguments"] = [ "sh", "-c", plist.command]
      if intervals = calendar_intervals
        attrs["StartCalendarInterval"] = calendar_intervals
      end
    end
  end

  def calendar_intervals
    crontab_expression = CrontabExpression.new(
      :minute => plist.minute,
      :hour => plist.hour,
      :day => plist.day_of_month,
      :weekday => plist.weekday_list,
      :month => plist.month_list
    )

    intervals = crontab_expression.intervals.map do |i|
      capitalized = i.map { |k,v| [k.capitalize, v] }
      Hash[capitalized]
    end

    if intervals.empty?
      nil
    elsif intervals.size == 1
      intervals.first
    else
      intervals
    end
  end
end
