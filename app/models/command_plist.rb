class CommandPlist

  # Initialize a CommandPlist
  #
  # command - a Command representing the command to serialize to plist
  def initialize(command)
    @command = command
  end

  # Generate a plist representation of the command
  #
  # Returns an XML string
  def plist
    {
      "Label" => label,
      "ProgramArguments" => program_arguments
    }.to_plist
  end

  protected

  # Returns an underscored label for a command
  def label
    "com.zerowidth.launched." + @command.name.parameterize("_")
  end

  # Returns the program arguments for the command
  def program_arguments
    ["sh", "-c", @command.command]
  end

end
