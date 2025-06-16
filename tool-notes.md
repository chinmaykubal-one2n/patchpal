patchpal/
│
├── cmd/                    # Entry point (main.go)
│
├── api/                    # HTTP handlers
│   └── fix.go              # Handler for /fix endpoint
│
├── service/                # Core logic
│   ├── scanner.go          # Trivy scan logic
│   ├── fixer.go            # LLM integration logic
│
├── utils/                  # Utility functions (file handling, logging, etc.)
│   └── file.go
│
├── config/                 # App config/env loading
│   └── config.go
│
├── go.mod
└── go.sum


POINTS TO REMEMBER (LIMITATIONS)
1. TAKING 1 FILE AS OF NOW
2. FILE AND ERROR HANDLING IS NOT OPTIMIZED (NO ERROR HANDLING FOR NOW AND FILE READING IS IN ONE SHOT I.E READING ALL THE CONTENT AT ONCE)

<!-- test file upload -->
curl -X POST -F 'file=@yourfile.yaml' http://localhost:8080/fix
curl -X POST -F 'file=@vulnerable-manifests/k8s-manifets.yml' http://localhost:8080/fix
