{{- /* Parameters */ -}}
{{- $data := . -}}
{{- $startUid := 5000 -}}
{{- $baseDn := "ou=Users,dc=mydom,dc=com" -}}
{{- /* Template */ -}}
{{- range $idx, $row := $data.CsvRows -}}
dn: uid={{ fromCsv "Username" $row }},{{ $baseDn }}
objectClass: account
objectClass: posixAccount
cn: {{ fromCsv "First Name" $row }} {{ fromCsv "Last Name" $row }}
givenName: {{ fromCsv "First Name" $row }}
sn: {{ fromCsv "Last Name" $row }}
uid: {{ fromCsv "Username" $row }}
uidNumber: {{ add $startUid $idx }}
gidNumber: 10000
homeDirectory: /home/{{ fromCsv "Username" $row }}
loginShell: /bin/bash
userPassword: {{ fromCsv "Password" $row | encodeLdifPassword }}

{{ end -}}
