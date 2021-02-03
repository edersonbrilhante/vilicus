# Vilicus

## Run locally
```bash
docker-compose -f local-dev/docker-compose.yaml up -d
```

## Example of analysis
```bash
 curl -XPOST 'http://localhost:8040/container-scanning/analysis' \
-H 'Content-Type: application/json' \
-d '{"image":"node"}'
```

<details>
  <summary>Example Result</summary>
  
  ```json
    {
      "id": "be89226e-ff60-4e04-8804-e091529742c3",
      "image": "node",
      "status": "finished",
      "created_at": "2021-02-02T20:02:20.775067Z",
      "updated_at": "2021-02-02T20:07:11.059549Z",
      "vilicus_results": {
        "clair": {
          "unknown_vulns": [{
            "fix": "0:0",
            "urls": [
              "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-0501"
            ],
            "name": "CVE-2018-0501",
            "severity": "Unknown",
            "package_name": "apt",
            "package_version": "1.4.11"
          }]
        },
        "anchore_engine ": {
          "high_vulns": [{
              "fix": "None",
              "urls": [
                "https://security-tracker.debian.org/tracker/CVE-2020-27843"
              ],
              "name": "CVE-2020-27843",
              "severity": "High",
              "package_name": "libopenjp2-7",
              "package_version": "2.1.2-1.1+deb9u5"
            }
          ]
        },
        "trivy": {
          "high_vulns": [{
              "fix": "",
              "urls": [
                "https://gcc.gnu.org/viewcvs/gcc/trunk/gcc/config/arm/arm-protos.h?revision=266379&view=markup"
              ],
              "name": "CVE-2018-12886",
              "severity": "High",
              "package_name": "cpp-6",
              "package_version": "6.3.0-18+deb9u1"
            }
          ]
        }
      }
    }
  ```
</details>

