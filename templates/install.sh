{{- $filename := printf "%s.plist" .Plist.Label -}}
{{- $url := printf "%s/plists/%s.xml" .RootURL .Plist.ID -}}
set -e
file='{{ $filename }}'
url='{{ $url }}'
dir="$HOME/Library/LaunchAgents"
echo "downloading $file..."
mkdir -p "$dir"
curl -fsS -o "$dir/$file" "$url"
echo "installing $file..."
launchctl load -w "$dir/$file"
