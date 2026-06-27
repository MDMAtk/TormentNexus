sed -i 's/\(skill\.IsRetired = true\)/\1\n			skill.Successes = 0\n			skill.Failures = 0/g' go/internal/skillregistry/evolution.go
