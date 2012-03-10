class CreateLaunchdPlists < ActiveRecord::Migration
  def change
    create_table :launchd_plists do |t|
      t.string :uuid,     :limit => 36,   :null => false
      t.string :command,  :limit => 1024, :null => false
      t.string :name,                     :null => false

      t.string :minute
      t.string :hour
      t.string :day_of_month
      t.string :weekdays
      t.string :months
      t.integer :interval

      t.boolean :run_at_load
      t.boolean :launch_only_once

      t.string :user
      t.string :group
      t.string :root_directory
      t.string :working_directory

      t.timestamps
    end
  end
end
