## Usage ##    
	
Loading dictionaries from files
```go	
	err = i18n.InitFromDir(`en`, `/usr/lib/app/translations`)
	if err != nil {
		log.Println(`Dictionary loading error`, err)
	}
```

JSON dictionary template
```json
{
  "errors": {
    "internal_error": "Server error, please try again later"
  },
  "errors.user.signup": {
    "username_to_short": "Username to short",
    "form_min_age": "Minimum age is {min}",
    "form_max_age": "Maximum age is {max}",
    "form_min_length": "{field} minimum length is {min}",
    "form_max_length": "{field} maximum length is {max}"
  }
}
```

Using dictionary
```go
package main

import (
	"github.com/censync/go-i18n"
	"log"
	"os"
)

const (
	defaultLanguage = "en_US"

	// Json dictionaries location
	translationsPath = "/usr/lib/my_app/translations"
)

func main() {
	err := i18n.InitFromDir(defaultLanguage, translationsPath, `en_US`, `cs_CZ`)
	if err != nil {
		log.Println(`Loading dictionaries error`, err)
		os.Exit(1)
	}
	userLang, _ := os.LookupEnv("LANG")

	// Get translator for locale
	tr := i18n.New(userLang)

	// Creating simple string
	//
	//	Sample dict data:
	//	{
	//		"form.signup" : {
	//			"welcome" : "Welcome to registration"
	//		}
	//	}
	strSimple := tr.T("form.signup", "welcome")

	// Creating formatted string
	//
	//	Sample dict data with arguments:
	//	{
	//		"form.login" : {
	//			"connections_limit" : "Hello, {name}"
	//		}
	//	}
	strFormatted := tr.Tf("form.login", "connection_lost", i18n.M{"{name}": "John"})

	// Creating simple error
	errSimple := tr.ErrT("errors.connections", "connection_lost")

	// Creating formatted error
	//
	//
	//	Sample dict data:
	//	{
	//		"errors.connections" : {
	//			"connections_limit" : "Connections limit is {count}"
	//		}
	//	}
	errFormatted := tr.ErrTf("errors.connections", "connections_limit", i18n.M{"{count}": 50})

	log.Println(strSimple, strFormatted, errSimple, errFormatted)
}

```



Using advanced errors with translator
```go
package main

import (
	"github.com/censync/go-i18n"
	"log"
	"net/http"
	"os"
)

const (
	defaultLanguage = "en_US"
)

func main() {
	collection := i18n.DictionaryCollection{
		"en_US": {
			"errors": {
				"unknown": "Unknown error",
			},
			"errors.connections": {
				"connections_limit": "Connections limit is {count}",
			},
			"form.signup": {
				"welcome": "Welcome to registration",
				"disabled": "Registration is temporarily unavailable",
			},
			"form.login": {
				"title": "Hello, {name}",
			},
		},
		"cs_CZ": {
			"errors": {
				"unknown": "Neznámá chyba",
			},
			"errors.connections": {
				"connections_limit": "Limit připojení je {count}",
			},
			"form.signup": {
				"welcome":  "Vítejte v registraci",
				"disabled": "Registrace je dočasně nedostupná",
			},
			"form.login": {
				"title": "Ahoj, {name}",
			},
		},
	}

	err := i18n.Init(defaultLanguage, collection)
	if err != nil {
		log.Println(`Loading dictionaries error`, err)
		os.Exit(1)
	}
	userLang, _ := os.LookupEnv("LANG")

	// Get translator for locale
	tr := i18n.New(userLang)

	errAdvanced := i18n.NewErr("form.signup", "disabled")

	log.Println(errAdvanced.T(tr)) // Returns string "Registration is disabled"

	log.Println(errAdvanced.ErrT(tr)) // Returns error "Registration is disabled"

	errAdvanced = i18n.NewErrWithCode(
		http.StatusTooManyRequests,
		"errors.connections",
		"connections_limit",
		i18n.M{"{count}": 50},
	)

	log.Println(errAdvanced.T(tr)) // Returns string "Connections limit is 50"

	log.Println(errAdvanced.ErrT(tr)) // Returns error "Connections limit is 50"
}

```