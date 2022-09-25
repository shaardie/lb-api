defaults
    mode tcp
    timeout client 10s
    timeout connect 5s
    timeout server 10s
{{- range .}}
    {{- $name := .Name }}
    {{- range $i, $f := .Config.Frontends }}
    {{- $port := .Port }}
frontend frontend-{{ $name }}-{{ $port }}
    bind :{{ $port }}
    default_backend backend-{{ $name }}-{{ $port }}
backend backend-{{ $name }}-{{ $port }}
        {{- range $j, $s := .Backend.Server }}
    server server{{ $j }} {{ $s }}
        {{- end }}
    {{- end }}
{{- end }}
