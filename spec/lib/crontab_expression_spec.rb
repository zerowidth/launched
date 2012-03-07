require "spec_helper"

require "crontab_expression"

describe CrontabExpression do

  def intervals(expr)
    CrontabExpression.new(expr).intervals
  end

  describe '#intervals' do
    it "returns a single minute" do
      intervals(:minute => "0").should == [{ :minute => 0 }]
    end

    it "parses a list of minutes" do
      intervals(:minute => "0,20,40").should == [
        {:minute => 0},
        {:minute => 20},
        {:minute => 40},
      ]
    end

    it "parses a divisor for minutes" do
      intervals(:minute => "*/5").should ==
        12.times.map { |n| {:minute => n * 5} }
    end

    it "parses a range of hours" do
      intervals(:hour => "12-15").should == [
        {:hour => 12},
        {:hour => 13},
        {:hour => 14},
        {:hour => 15}
      ]
    end

    it "handles a divisor for a range" do
      intervals(:hour => "12-16/2").should == [
        {:hour => 12},
        {:hour => 14},
        {:hour => 16}
      ]
    end

    it "handles the same minute for two different hours" do
      intervals(:minute => "30", :hour => "0,12").should == [
        {:minute => 30, :hour => 0},
        {:minute => 30, :hour => 12}
      ]
    end

    it "returns the cartesian product of each specified minute, hour, and day" do
      intervals(:minute => "0,30", :hour => "12-13", :day => "1,15").should == [
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
      intervals(:minute => "0", :hour => "*").should == [
        {:minute => 0}
      ]
    end

    it "ignores '*' in compound expressions, since it overrides all" do
      intervals(:minute => "0", :hour => "5,*,20").should == [
        {:minute => 0}
      ]
    end
  end

end
