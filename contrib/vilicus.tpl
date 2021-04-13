{
  "version": "14.0.0",
  "vulnerabilities": [
{{- $first_rule := true }}
{{- range $key, $analysis := . }}
  {{- $vulns := vulnList $analysis.Results  }}
  {{- range $vuln_index, $vuln := $vulns }}
    {{- if $first_rule -}}
      {{- $first_rule = false -}}
    {{ else -}}
      ,
    {{- end }}
    {
      "vendor": "{{ $vuln.Vendor }}",
      "severity": "{{ $vuln.Severity }}",
      "name": "{{ $vuln.Name }}",
      "package_name": "{{ $vuln.PackageName }}",
      "package_version": "{{ $vuln.PackageVersion }}",
      "solution": {{ if $vuln.Fix -}}
                    "Upgrade {{ $vuln.PackageName }} to {{ $vuln.Fix }}"
                  {{- else -}}
                    "No solution provided"
                  {{- end }}
    }
  {{- end }}
{{- end }}
   ]
}