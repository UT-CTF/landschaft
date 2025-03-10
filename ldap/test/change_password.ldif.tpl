{{- /* Parameters */ -}}
{{- $data := . -}}
{{- $baseDn := "ou=Users,dc=mydom,dc=com" -}}
{{- /* Template */ -}}
{{- range $_, $row := $data.CsvRows -}}
dn: uid={{ fromCsv "Username" $row }},{{ $baseDn }}
changetype: modify
replace: userPassword
userPassword: {{ fromCsv "Password" $row | encodeLdifPassword }}

{{ end -}}
