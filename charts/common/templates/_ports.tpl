
{{- /*
common.ports takes a list of dicts and turns them into ports.
*/ -}}
{{- define "common.ports" -}}
{{ if . }}
ports:
  -
{{ range $k, $v := . }}
{{ range $kk, $vv := $v }}
   {{ $kk }}: {{ $vv }}
{{- end }}
{{- end }}
{{- end }}
{{- end -}}