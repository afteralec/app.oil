# Oil App

```bash
curl -LJO https://raw.githubusercontent.com/afteralec/db.oil/main/schema.hcl \
&& atlas migrate diff --dev-url "mysql://root:pass@:3306/clean" --to "file://schema.hcl"
```
