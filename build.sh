find ./publisher/**/* -type f -name '*.go' -print0 | while IFS= read -r -d '' file; do
  echo "build plugin $file" 
  go build -o ./plugins -buildmode=plugin $file
done