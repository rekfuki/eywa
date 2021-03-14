
{{- /*
common.ports takes a list of dicts and turns them into ports.
*/ -}}
{{- define "common.ports" -}}
ports:
  -
{{ range $k, $v := . }}
{{ range $kk, $vv := $v }}
   {{ $kk }}: {{ $vv }}
{{- end }}
{{- end }}
{{- end -}}