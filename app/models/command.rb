class Command
  include ActiveModel::Validations

  attr_accessor :name, :command

  validates_presence_of :name, :command

  def initialize(params)
    @name = params[:name]
    @command = params[:command]
  end

end
