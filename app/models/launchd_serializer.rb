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

      if plist.interval
        attrs["StartInterval"] = plist.interval
      end

      if plist.run_at_load
        attrs["RunAtLoad"] = true
      end

      if plist.launch_only_once
        attrs["LaunchOnlyOnce"] = true
      end
    end
  end

  def calendar_intervals
    crontab_expression = CrontabExpression.new(
      :minute => plist.minute,
      :hour => plist.hour,
      :day => plist.day_of_month,
      :day_of_week => plist.weekday_list,
      :month => plist.month_list
    )

    intervals = crontab_expression.intervals.map do |i|
      capitalized = i.map do |k,v|
        if k.to_s == "day_of_week"
          ["Weekday", v]
        else
          [k.capitalize, v]
        end
      end
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
