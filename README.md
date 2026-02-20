# go-i18n

Lightweight internationalization (i18n) library for Go applications. Provides translation management
with support for multiple locales, parameterized messages, structured i18n-aware errors, and JSON serialization.

## Features

- Dictionary-based translations organized by sections and keys
- Loading translations from JSON files or in-memory maps
- Parameterized message formatting with type-aware placeholder substitution
- Structured i18n error types (`I18nError`, `I18nMultipleError`) with JSON marshaling
- Locale-aware JSON serialization of errors
- Builder pattern for error construction
- Concurrent-safe access to translation data
- Protobuf schema for cross-service error transport

## Installation

```
go get github.com/censync/go-i18n
```

## Quick Start

### Initialization from JSON Files

Place JSON dictionary files in a directory. Each file name corresponds to a locale
(e.g., `en_US.json`, `cs_CZ.json`).

```go
package main

import (
	"log"
	"os"

	"github.com/censync/go-i18n"
)

func main() {
	// Load all .json files from the directory as locale dictionaries.
	// The first argument is the default (fallback) locale.
	err := i18n.InitFromDir("en_US", "/usr/lib/my_app/translations")
	if err != nil {
		log.Println("Loading dictionaries error:", err)
		os.Exit(1)
	}
}
```

You can also specify an explicit list of locales to load:

```go
err := i18n.InitFromDir("en_US", "/usr/lib/my_app/translations", "en_US", "cs_CZ")
```

### Initialization from Code

```go
collection := i18n.DictionaryCollection{
	"en": &i18n.Dictionary{
		"errors": &i18n.DictionaryEntry{
			"unknown": "Unknown error",
		},
		"form.signup": &i18n.DictionaryEntry{
			"welcome": "Welcome to registration",
		},
	},
	"cs": &i18n.Dictionary{
		"errors": &i18n.DictionaryEntry{
			"unknown": "Neznama chyba",
		},
		"form.signup": &i18n.DictionaryEntry{
			"welcome": "Vitejte v registraci",
		},
	},
}

err := i18n.Init("en", &collection)
```

### JSON Dictionary Format

```json
{
  "errors": {
    "internal_error": "Server error, please try again later"
  },
  "errors.user.signup": {
    "username_too_short": "Username too short",
    "form_min_age": "Minimum age is {min}",
    "form_max_age": "Maximum age is {max}",
    "form_min_length": "{field} minimum length is {min}",
    "form_max_length": "{field} maximum length is {max}"
  }
}
```

## Usage

### Basic Translation

```go
// Get a translator for the user locale.
// Falls back to the default locale if the requested one is not available.
tr := i18n.Get("cs_CZ")

// Simple key lookup: returns the translated string,
// or "section.key" if not found.
msg := tr.T("form.signup", "welcome")
```

### Formatted Translation

Use placeholders wrapped in braces inside dictionary values. Pass an `i18n.M` map
to substitute them at runtime. Integer, float, and string values are supported.

```go
// Dictionary value: "Hello, {name}! You have {count} messages."
msg := tr.Tf("form.login", "greeting", i18n.M{
	"{name}":  "John",
	"{count}": 5,
})
// Result: "Hello, John! You have 5 messages."
```

### Translated Errors

```go
// Simple translated error
err := tr.ErrT("errors", "unknown")
// err.Error() == "Unknown error"

// Formatted translated error
err = tr.ErrTf("errors.connections", "connections_limit", i18n.M{"{count}": 50})
// err.Error() == "Connections limit is 50"
```

### Available and Default Locales

```go
locales := i18n.AvailableLocales() // []string{"en", "cs"}
def := i18n.DefaultLocale()        // "en"
```

## Structured Errors

The library provides `I18nError` and `I18nMultipleError` types designed for API responses.
They carry section/key references and can be serialized to JSON either as raw keys
or as translated messages (when a locale is set).

### I18nError

```go
// Create a simple error
e := i18n.NewErr("form.signup", "disabled")
e.Error() // "form.signup.disabled"

// Create an error with HTTP status code and format values
e = i18n.NewErrWithCode(http.StatusTooManyRequests, "errors.connections", "connections_limit",
	i18n.M{"{count}": 50},
)
e.Code() // 429

// Builder pattern
e = i18n.NewErr("errors", "unknown").
	WithCode(http.StatusInternalServerError).
	WithLocale("cs_CZ")
```

### JSON Serialization

Without a locale set, errors serialize as `{"section":"key"}`:

```go
e := i18n.NewErr("form.signup", "disabled")
b, _ := json.Marshal(e)
// b == {"form.signup":"disabled"}
```

With a locale set, errors serialize as the translated string:

```go
e := i18n.NewErr("form.signup", "disabled").WithLocale("en")
b, _ := json.Marshal(e)
// b == "Registration is temporarily unavailable"
```

### Translating Errors Explicitly

```go
tr := i18n.Get("cs_CZ")
e := i18n.NewErr("form.signup", "welcome")

msg := e.T(tr)     // translated string
err := e.ErrT(tr)  // translated error

e2 := i18n.NewErrWithCode(429, "errors.connections", "connections_limit",
	i18n.M{"{count}": 50},
)
msg = e2.Tf(tr)     // formatted translated string
err = e2.ErrTf(tr)  // formatted translated error
```

### I18nMultipleError

Represents multiple field-level errors, useful for form validation responses.

```go
// Create with a single field error
me := i18n.NewMultipleErr("username", "fields.errors", "too_short")

// Or start empty and add errors
me = i18n.NewMultipleEmptyErr()
me.Add("username", "fields.errors", "too_short")
me.Add("email", "fields.errors", "invalid_format")
me.AddDefault("errors", "validation_failed")

if me.HasErrors() {
	b, _ := json.Marshal(me)
	// Without locale: {"username":{"fields.errors":"too_short"},"email":{"fields.errors":"invalid_format"}, ...}
}
```

With a locale, each field error is translated:

```go
me := i18n.NewMultipleErr("username", "fields.errors", "too_short")
me.SetLocale("en")
b, _ := json.Marshal(me)
// {"username":"Field too short"}
```

## Protobuf Support

The `i18n.proto` file provides Protocol Buffers definitions for `I18nError` and
`I18nMultipleError`, enabling transport of structured i18n errors across gRPC services.

## Thread Safety

All public functions and methods are safe for concurrent use. The library uses
`sync.RWMutex` internally to protect shared translation data. Initialization
functions (`Init`, `InitFromDir`) acquire a write lock; all read operations
(`Get`, `T`, `Tf`, `AvailableLocales`, `DefaultLocale`) use read locks.

## License

MIT License. See [LICENSE](LICENSE) for details.