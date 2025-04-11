{{- /* Parameters */ -}}
{{- $data := . -}}
{{- $baseDn := "ou=Users,dc=mydom,dc=com" -}}
{{- /* Template */ -}}
dn: uid={{ $data.Username }},{{ $baseDn }}
changetype: modify
replace: userPassword
userPassword: {{ $data.Password | encodeLdifPassword }}
