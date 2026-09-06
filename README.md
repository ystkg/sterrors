# スタックトレース付きエラー

## 概要

Go言語の `error` にスタックトレースを付けるスケルトンコードです。
個別の用途に応じてカスタマイズして使うことを想定していますが、そのままで使うことも可能です。

## 推奨環境

Go version 1.27以降

例

```ShellSession
$ go version
go version go1.27.1 linux/amd64
```

## 使用例

`"github.com/ystkg/sterrors"`をimportします。

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/ystkg/sterrors"
)

func ParseInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, sterrors.WithFrames(err)
	}
	return int16(i), nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "stackTraces" {
				if err, ok := a.Value.Any().(error); ok {
					errs := sterrors.StackTraces(err)
					return slog.Any(a.Key, sterrors.Format(errs))
				}
			}
			return a
		},
	})))

	ctx := context.Background()

	_, err := ParseInt16("abc")
	if err != nil {
		slog.ErrorContext(ctx, "ParseInt16", "stackTraces", err)
	}
}
```

### 実行結果の例

```json
{
  "time": "2026-09-06T09:00:02.542217674+09:00",
  "level": "ERROR",
  "msg": "ParseInt16",
  "stackTraces": [
    {
      "error": "strconv.ParseInt: parsing \"abc\": invalid syntax",
      "stackTrace": [
        "main.ParseInt16(ex1/main.go:15)",
        "main.main(ex1/main.go:35)",
        "runtime.main(runtime/proc.go:302)",
        "runtime.goexit(runtime/asm_amd64.s:1264)"
      ]
    }
  ]
}
```
