{{- /* Parameters */ -}}
{{- $data := . -}}
{{- $groupName := "employees" -}}
{{- $baseDn := "ou=Groups,dc=mydom,dc=com" -}}
{{- /* Template */ -}}
{{- range $_, $row := $data.CsvRows -}}
{{- $csvGroups := fromCsv "Groups" $row -}}
{{- if contains $groupName $csvGroups -}}
dn: cn={{ $groupName }},{{ $baseDn }}
changetype: modify
add: memberUid
memberUid: {{ fromCsv "User" $row  }}

{{ end -}}
{{- end -}}
