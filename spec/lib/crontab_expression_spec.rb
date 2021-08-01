require "spec_helper"

describe CrontabExpression do

  def intervals(expr)
    CrontabExpression.new(expr).intervals
  end

  describe '#intervals' do
    it "returns a single minute" do
      expect(intervals(:minute => "0")).to match_array([{ :minute => 0 }])
    end

    it "parses a list of minutes" do
      expect(intervals(:minute => "0,20,40")).to match_array([
        {:minute => 0},
        {:minute => 20},
        {:minute => 40},
      ])
    end

    it "parses a divisor for minutes" do
      expect(intervals(:minute => "*/5")).to match_array(
        12.times.map { |n| {:minute => n * 5} }
      )
    end

    it "parses a range of hours" do
      expect(intervals(:hour => "12-15")).to match_array [
        {:hour => 12},
        {:hour => 13},
        {:hour => 14},
        {:hour => 15}
      ]
    end

    it "handles a divisor for a range" do
      expect(intervals(:hour => "12-16/2")).to match_array [
        {:hour => 12},
        {:hour => 14},
        {:hour => 16}
      ]
    end

    it "handles the same minute for two different hours" do
      expect(intervals(:minute => "30", :hour => "0,12")).to match_array [
        {:minute => 30, :hour => 0},
        {:minute => 30, :hour => 12}
      ]
    end

    it "returns the cartesian product of each specified minute, hour, and day" do
      expect(intervals(:minute => "0,30", :hour => "12-13", :day => "1,15")).to match_array [
        {:minute => 0, :hour => 12, :day => 1},
        {:minute => 0, :hour => 12, :day => 15},
        {:minute => 0, :hour => 13, :day => 1},
        {:minute => 0, :hour => 13, :day => 15},
        {:minute => 30, :hour => 12, :day => 1},
        {:minute => 30, :hour => 12, :day => 15},
        {:minute => 30, :hour => 13, :day => 1},
        {:minute => 30, :hour => 13, :day => 15},
      ]
    end

    it "ignores '*' in expressions" do
      expect(intervals(:minute => "0", :hour => "*")).to match_array [
        {:minute => 0}
      ]
    end

    it "ignores '*' in compound expressions, since it overrides all" do
      expect(intervals(:minute => "0", :hour => "5,*,20")).to match_array [
        {:minute => 0}
      ]
    end
  end

end
