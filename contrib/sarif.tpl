{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
  {{- $first_analysis := true }}
  {{- range $key, $analysis := . }}
    {{- if $first_analysis -}}
      {{- $first_analysis = false -}}
    {{ else -}}
      ,
    {{- end }}
    {
      "tool": {
        "driver": {
          "name": "Vilicus",
          "informationUri": "https://github.com/edersonbrilhante/vilicus",
          "fullName": "Vilicus",
          "semanticVersion": "v0.0.1", 
          "version": "0.0.1", 
          {{- $vulns := vulnList $analysis.Results  }}
          "rules": [
          {{- $first_rule := true }}
          {{- range $vuln_index, $vuln := $vulns }}
            {{- if $first_rule -}}
              {{- $first_rule = false -}}
            {{ else -}}
              ,
            {{- end }}
            {{- $help_link := ""}}
            {{ $length := len $vuln.URL }} {{ if gt $length 0 }}
              {{- $help_link = index $vuln.URL 0 -}}
            {{- end }}
            {
              "id": "[VULICUS][{{ $vuln.Vendor }}][{{ $vuln.Severity }}] {{ $vuln.Name }} {{ $vuln.PackageName }} {{ $vuln.PackageVersion }}",
              "name": "dockerfile_scan",
              "shortDescription": {
                "text": {{ printf "%v %v Severity: %v Package: %v Version: %v" $vuln.Vendor $vuln.Name $vuln.Severity $vuln.PackageName $vuln.PackageVersion | printf "%q" }}
              },
              "fullDescription": {
                "text": {{ printf "%v %v Severity: %v Package: %v Version: %v" $vuln.Vendor $vuln.Name $vuln.Severity $vuln.PackageName $vuln.PackageVersion | printf "%q" }}
              },
              "help": {
                "text": {{ printf "Vulnerability: %v\nSeverity: %v\nPackage: %v\nVersion: %v\nFixed: %v\nLink: [%v](%v)\n" $vuln.Name $vuln.Severity $vuln.PackageName $vuln.PackageVersion $vuln.Fix $vuln.Name $help_link | printf "%q" }},
                "markdown": {{ printf "**Vulnerability %v**\n| Severity | Package | Version | Fixed | Link |\n| --- | --- | --- | --- | --- |\n| %v | %v | %v | %v | [%v](%v) |\n" $vuln.Name $vuln.Severity $vuln.PackageName $vuln.PackageVersion $vuln.Fix $vuln.Name $help_link | printf "%q" }}
              },
              "properties": {
                "tags": [
                  "vulnerability",
                  "{{ $vuln.Severity }}",
                  "{{ $vuln.Vendor }}",
                  "{{ $vuln.Name }}",
                  {{ $vuln.PackageName | printf "%q" }}
                ],
                "precision": "very-high"
              }
            }
          {{- end }}
          ]
        }
      },
      "automationDetails": {
        "description": {
          "text": "This is the run  {{ $analysis.Image }}"
        },
        "id": "{{ $analysis.ID }}",
        "guid": "{{ $analysis.ID }}",
        "properties": {
          "tags": [
            "vulnerability",
            "vilicus",
            "{{ $analysis.Image }}"
          ]
        }
      },
      "results":[
      {{- $first_rule := true }}
      {{- range $vuln_index, $vuln := $vulns }}
        {{- if $first_rule -}}
          {{- $first_rule = false -}}
        {{ else -}}
          ,
        {{- end }}
        {
          "ruleId": "[VULICUS][{{ $vuln.Vendor }}][{{ $vuln.Severity }}] {{ $vuln.Name }} {{ $vuln.PackageName }} {{ $vuln.PackageVersion }}",
          "ruleIndex": {{ $vuln_index }},
          "level": "error",
          "kind": "fail",
          "message": {
            "text": {{ printf "%v - %v Package: %v Version: %v" $vuln.Vendor $vuln.Name $vuln.PackageName $vuln.PackageVersion | printf "%q" }}
          },
          "locations": [{
            "physicalLocation": {
              "artifactLocation": {
                "uri": "Dockerfile"
              },
              "region": {
                "startLine": 1,
                "startColumn": 1,
                "endColumn": 1
              }
            },
            "logicalLocations": [
              {
                "fullyQualifiedName": "dockerfile"
              }
            ]
          }]
        }
      {{- end }}
      ],
      "columnKind": "utf16CodeUnits"
    }
  {{- end }}
  ]
}