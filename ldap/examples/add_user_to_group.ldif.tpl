{{ $data := . }}
{{- range $_, $row := $data.CsvRows -}}
{{- if fromCsv "Groups" $row | contains "employees" -}}
dn: cn=employees,ou=Groups,dc=mydom,dc=com
changetype: modify
add: memberUid
memberUid: {{ fromCsv "User" $row  }}
{{- end -}}
{{ end -}}
