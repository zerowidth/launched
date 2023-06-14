{{- $filename := printf "%s.plist" .Plist.Label -}}
{{- $url := printf "%s/plist/%s.xml" .RootURL .Plist.Encode -}}
echo downloading {{ $filename }}...
mkdir -p ~/Library/LaunchAgents
curl -o ~/Library/LaunchAgents/{{ $filename }} {{ $url }}
echo installing {{ $filename }}...
launchctl load -w ~/Library/LaunchAgents/{{ $filename }}
