#!/bin/bash
terraform output -json > output_temp.json
echo "{" > formatted_output.json
grep -Po '"[^"]+":\s*\{\s*"value":\s*"[^"]+"' output_temp.json | while read -r line; do
  key=$(echo "$line" | grep -Po '"[^"]+"' | head -1 | tr -d '"')
  value=$(echo "$line" | grep -Po '"value":\s*"[^"]+"' | sed 's/"value":\s*"\(.*\)"/\1/')
  if [ -n "$value" ]; then
    echo "  \"$key\": \"$value\"," >> formatted_output.json
  fi
done
sed -i '$ s/,$//' formatted_output.json
echo "}" >> formatted_output.json
cat formatted_output.json
rm output_temp.json
