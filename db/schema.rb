# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2012_03_10_202133) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "launchd_plists", id: :serial, force: :cascade do |t|
    t.string "uuid", limit: 36, null: false
    t.string "command", limit: 1024, null: false
    t.string "name", null: false
    t.string "minute"
    t.string "hour"
    t.string "day_of_month"
    t.string "weekdays"
    t.string "months"
    t.integer "interval"
    t.boolean "run_at_load"
    t.boolean "launch_only_once"
    t.string "user"
    t.string "group"
    t.string "root_directory"
    t.string "working_directory"
    t.datetime "created_at"
    t.datetime "updated_at"
  end

end
