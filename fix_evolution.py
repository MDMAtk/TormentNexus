with open("go/internal/skillregistry/decision.go", "r") as f:
    content = f.read()

content = content.replace("	IsRetired  bool       `json:\"useCount\"`", "	UseCount   int        `json:\"useCount\"`\n	IsRetired  bool       `json:\"isRetired\"`")
content = content.replace("	UseCount   int\n", "")

with open("go/internal/skillregistry/decision.go", "w") as f:
    f.write(content)
