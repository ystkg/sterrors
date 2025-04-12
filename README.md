# スタックトレース付きエラー

## 概要

Go言語の `error` にスタックトレースを付けるスケルトンコードです。
個別の用途に応じてカスタマイズして使うことを想定していますが、そのままで使うことも可能です。

## 推奨環境

Go version 1.24以降

例

```ShellSession
$ go version
go version go1.24.2 linux/amd64
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
  "time": "2025-04-12T09:37:07.55737061+09:00",
  "level": "ERROR",
  "msg": "ParseInt16",
  "stackTraces": [
    {
      "error": "strconv.ParseInt: parsing \"abc\": invalid syntax",
      "stackTrace": [
        "main.ParseInt16(/main/main.go:15)",
        "main.main(/main/main.go:35)",
        "runtime.main(/usr/local/go/src/runtime/proc.go:283)",
        "runtime.goexit(/usr/local/go/src/runtime/asm_amd64.s:1700)"
      ]
    }
  ]
}
```
