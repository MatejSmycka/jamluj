# JAMLUJ

This tool is CLI Runner for jobs, which are defined in YAML file.

# Example

```yaml
---
run: # List of jobs to run in sequence
  - command:
      val: echo 1
  - command:
      val: # List of commands to run in parallel
        - echo 2
        - sleep 10
        - echo 3
  - command:
      val: echo4.py 
      python: true # Run python script, command is prefixed with `poetry run python`
```
