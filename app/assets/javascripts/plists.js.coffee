log = (args...) ->
  if console.log
    console.log args...

cronValidator = (value, element, params...) ->
  value is "" or value.match ///
    ^ \s*
    # define a basic cron "atom" - *, */5, 1-10, 1-10/2, etc.
    ( \* (/\d+)? | \d+ (-\d+)? (/\d+)? )
    # followed by zero or more of the same, separated by commas
    ( , ( \* (/\d+)? | \d+ (-\d+)? (/\d+)? ) )*
    \s* $
  ///

cronRange = (value, element, [min, max]) ->
  return true unless value
  digits = (parseInt(n, 10) for n in value.match /\d+/g)
  if digits
    rejected = (d for d in digits when d < min or d > max)
    rejected.length is 0
  else
    true

$.validator.addMethod "cron", cronValidator, "Must be a cron expression"
$.validator.addMethod "cron_range", cronRange, $.format "Digits must be between {0} and {1}"

$ ->
  $('#plist').button()
  $('#plist #label').focus().select()

  $('#plist .btn-group .btn').click (event) ->
    buttons = $(this).parent('.btn-group').children('.btn')

    # when it's clicked, it hasn't toggled yet. so:
    if $(this).hasClass 'active' # it's being toggled off
      selected = (i for element, i in buttons when $(element).hasClass("active") and element isnt this)
    else
      selected = (i for element, i in buttons when $(element).hasClass("active") or element is this)

    $(this).closest('.controls').find('input').val selected.join(",")

    event.preventDefault()


  $('#plist').validate
    errorElement: "span"
    errorClass: "error"
    successClass: "success"

    onkeyup: false # this is a little aggressive

    errorPlacement: (error, element) ->
      $err = $(error).addClass('help-inline')
      $(element).after $err

    highlight: (element, errorClass) ->
      $(element).closest('.control-group').removeClass('success').addClass('error')

    unhighlight: (element, successClass) ->
      group = $(element).closest('.control-group')
      if group.hasClass('error')
        group.removeClass('error').addClass('success')

    rules:
      label: "required"
      command: "required"
      interval: "digits"
      minute:
        cron: true
        cron_range: [0, 59]
      hour:
        cron: true
        cron_range: [0, 23]
      day_of_month:
        cron: true
        cron_range: [1, 31]

