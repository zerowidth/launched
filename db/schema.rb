# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20120310202133) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "launchd_plists", force: :cascade do |t|
    t.string   "uuid",              limit: 36,   null: false
    t.string   "command",           limit: 1024, null: false
    t.string   "name",                           null: false
    t.string   "minute"
    t.string   "hour"
    t.string   "day_of_month"
    t.string   "weekdays"
    t.string   "months"
    t.integer  "interval"
    t.boolean  "run_at_load"
    t.boolean  "launch_only_once"
    t.string   "user"
    t.string   "group"
    t.string   "root_directory"
    t.string   "working_directory"
    t.datetime "created_at",                     null: false
    t.datetime "updated_at",                     null: false
  end

end
