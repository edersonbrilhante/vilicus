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
      "id": "[VULICUS][{{ $vuln.Vendor }}][{{ $vuln.Severity }}] {{ $vuln.Name }} {{ $vuln.PackageName }} {{ $vuln.PackageVersion }}",
      "category": "container_scanning",
      "message": "[VULICUS][{{ $vuln.Vendor }}][{{ $vuln.Severity }}] {{ $vuln.Name }} {{ $vuln.PackageName }} {{ $vuln.PackageVersion }}",
      "description": {{ printf "[%v] %v Severity: %v Package: %v Version: %v" $vuln.Vendor $vuln.Name $vuln.Severity $vuln.PackageName $vuln.PackageVersion | printf "%q" }},
      "cve": "{{ $vuln.Name }}",
      "severity": {{ if eq $vuln.Severity "Negligible" -}}
                    "Info"                   
                  {{-  else -}}
                    "{{ $vuln.Severity }}"
                  {{- end }},
      "confidence": "High",
      "solution": {{ if $vuln.Fix -}}
                    "Upgrade {{ $vuln.PackageName }} to {{ $vuln.Fix }}"
                  {{- else -}}
                    "No solution provided"
                  {{- end }},
      "scanner": {
        "id": "{{ $vuln.Vendor }}",
        "name": "{{ $vuln.Vendor }}",
        "version": "",
        "vendor": {
          "name": "Vilicus"
        }
      },
      "location": {
        "dependency": {
          "package": {
            "name": "{{ $vuln.PackageName }}"
          },
          "version": "{{ $vuln.PackageVersion }}"
        },
        "operating_system": "Unknown",
        "image": "{{ $analysis.Image }}"
      },
      "identifiers": [
        {{- $help_link := ""}}
        {{ $length := len $vuln.URL }} {{ if gt $length 0 }}
            {{- $help_link = index $vuln.URL 0 -}}
        {{- end }}
        {
          "type": "cve",
          "name": "{{ $vuln.Name }}",
          "value": "{{ $vuln.Name }}",
          "url": "{{ $help_link }}"
        }
      ],
      "links": [
      {{- $first_url := true -}}
      {{- range $url_index, $url := $vuln.URL }}
        {{- if $first_url -}}
          {{- $first_url = false }}
        {{- else -}}
          ,
        {{- end -}}
        {
          "url": "{{ $url }}"
        }
      {{- end }}
      ]
    }
  {{- end }}
{{- end }}
  ],
  "remediations": [
  ]
}