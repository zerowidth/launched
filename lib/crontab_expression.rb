class CrontabExpression

  # the ranges for each cron granularity
  RANGES = {
    :minute => 0..59,
    :hour => 0..23,
    :day => 1..31,
    :month => 1..12,
    :day_of_week => 0..6
  }

  # The expression hash
  attr_reader :expression

  # Initialize a new crontab expression
  #
  # expression - a hash of cron-style columns, optionally including any or
  #              all of these, as cron expression strings:
  #                :minute
  #                :hour
  #                :day (of month)
  #                :month
  #                :day_of_week
  #
  # Blank or nil columns are left out of subsequent calculations.
  #
  # Examples:
  #
  #   {:minute => "0", :hour => "0"}
  #   {:minute => "*/20", :hour => "9-16/2"}
  #
  def initialize(expression={})
    @expression = expression
  end

  # Public: calculate the distinct times for the expression
  #
  # "*" entries are left out, as these are implicitly defined in the target
  # system (plist)
  #
  # Returns an array of unique minute/hour/day/month/day_of_week required
  # to implement this expression.
  def intervals
    parts = {}
    [:minute, :hour, :day, :month, :day_of_week].each do |granularity|
      if times = times_for(granularity)
        parts[granularity] = times
      end
    end

    merged_product parts
    # interval_list.concat times.map {|t| {granularity => t} }
    # m = expression[:minute]
    # m.split(",").each do |part|
    #   time, divisor = part.split("/")
    #   start_time, end_time = time.split "-"

    #   interval_list << {:minute => start_time.to_i}
    # end
  end

  # Calculate the values for the given cron expression
  #
  # granularity - what granularity this expression applies to, e.g. :minute,
  #               :hour, :day, :month, :day_of_week
  #
  # Returns an array of integer times, or nil if no expressions apply.
  def times_for(granularity)
    expr = expression[granularity]
    return nil unless expr
    return nil if expr == ""
    return nil if expr.split(",").include? "*"

    [].tap do |times|
      expr.split(",").each do |value|
        time, divisor = value.split "/"
        divisor = divisor.to_i if divisor
        start_time, end_time = time.split "-"
        if time == "*" # will always have a divisor
          divided = RANGES[granularity].select do |g|
            g % divisor == 0
          end
          times.concat divided
        else
          if end_time
            d = divisor || 1
            start_time.to_i.upto(end_time.to_i) { |t| times << t if t % d == 0 }
          else
            times << start_time.to_i
          end
        end
      end
    end
  end

  # Return the cartesian product of the given parts as merged hashes.
  #
  # parts - a hash of the time parts to combine
  #
  # Example:
  #
  #   cartesian_product(:m => [0,30], :h => [8,16]) # => [
  #     {:m => 0, :h => 8}
  #     {:m => 0, :h => 16}
  #     {:m => 30, :h => 8}
  #     {:m => 30, :h => 16}
  #   ]
  #
  # Returns an array of the merged/combined parts
  def merged_product(parts)
    result = []
    parts.map do |key, values|
      if result.empty?
        values.each { |v| result << {key => v} }
      else
        combined = []
        result.each do |r|
          values.each do |v|
            combined << r.merge(key => v)
          end
        end
        result = combined
      end
    end
    result
  end

end
