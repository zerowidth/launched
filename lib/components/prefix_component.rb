module PrefixComponent
  def prefix(wrapper_options = nil)
    @prefix ||= begin
      options[:prefix].to_s.html_safe if options[:prefix].present?
    end
  end
end

SimpleForm.include_component(PrefixComponent)
