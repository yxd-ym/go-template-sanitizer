# go-template-sanitizer

## Install

```
go install github.com/yxd-ym/go-template-sanitizer@latest
```

## Usage

The binary `go-template-sanitizer` reads the input from stdin and output the processed template to stdout.

Run

```
cat your-template-file.yaml | go-template-sanitizer | tee your-template-file-sanitized.yaml
```

To see how your template file is being sanitized.
