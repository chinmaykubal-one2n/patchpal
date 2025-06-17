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
curl -X POST -F 'file=@k8s-manifets.yml' http://localhost:8080/fix


IMPORTANT STEP AFTER GETTING THE OUTPUT FORM THE LLM BEFORE SAVING TO FILE MAKE SURE TO PASS THAT TOO TRIVY AGAIN FOR ANY VULNARABILITS (send again :- llm response yaml, current yaml and new trivy scan ).


THINGS WE CAN DO:-
🧪 Pro Tip: Add Diff Check in Your Code

After you get the fixed YAML from the LLM, do a programmatic diff (e.g. using sigs.k8s.io/yaml to parse and compare objects), and reject or retry if critical fields like resources or image are missing.


<!--  -->

