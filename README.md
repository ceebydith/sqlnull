# sqlnull

A Go package that provides a convenient way to handle SQL null values for various data types. This package simplifies the process of scanning SQL results into Go structs by wrapping your target variables and providing custom SQL scanners.

# Why?

Sometimes I get tired of special structs to handle nullable values ​​including `sql.NullString` and friends. I just need `*string` on my variables or structs, like mongodb can handle it very nice.

This package basicly based on `sql.NullString` and friends.

# Table of Contents
- [sqlnull](#sqlnull)
- [Why?](#why)
- [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [License](#license)
  - [Acknowledgements](#acknowledgements)

## Features
- **Supports various data types**: Including `bool`, `uint8`, `int8`, `int16`, `uint16`, `int32`, `uint32`, `int64`, `uint64`, `int`, `uint`, `string`, `float32`, `float64`, and `time.Time`.
- **Automatic zero values**: Sets target variables to their zero values if the SQL result is null.
- **Easy integration**: Simple to use with existing Go applications.
- **No dependency package**: Only use Go build-in package, except for testing, it use [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) and [`github.com/stretchr/testify`](https://github.com/stretchr/testify)

## Installation
1. Get the package
   ```sh
   go get -u github.com/ceebydith/sqlnull
   ```
2. Import the package
   ```go
   import "github.com/ceebydith/sqlnull"
   ```

## Usage
Here's how to use the `sqlnull` package in your Go application:
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/ceebydith/sqlnull"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Example usage with SQLite3
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    // Create users table with id and username as NOT NULL, phone and verified_at as NULL
    _, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY NOT NULL,
        username TEXT NOT NULL,
        phone TEXT,
        verified_at DATETIME
    )`)
    if err != nil {
        log.Fatal(err)
    }

    // Insert a sample user for demonstration purposes
    _, err = db.Exec(`INSERT INTO users (id, username, phone, verified_at) VALUES (1, 'johndoe', '123456789', NULL)`)
    if err != nil {
        log.Fatal(err)
    }

    type Customer struct {
        ID         int64
        Username   string
        Phone      *string
        VerifiedAt *time.Time
    }

    var cust Customer
    row := db.QueryRow("SELECT id, username, phone, verified_at FROM users WHERE id=?", 1)
    // for individual target use like below
    // err = row.Scan(sqlnull.Target(&cust.ID), sqlnull.Target(&cust.Username), sqlnull.Target(&cust.Phone), sqlnull.Target(&cust.VerifiedAt))
    err = row.Scan(sqlnull.Scanner(&cust.ID, &cust.Username, &cust.Phone, &cust.VerifiedAt)...)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Customer: %+v\n", cust)
}
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/ceebydith/sqlnull/blob/master/LICENSE) file for details.

## Acknowledgements
Special thanks to the Go community and contributors who made this project possible.


Feel free to customize the content to better fit your project's specific details, such as replacing placeholder URLs and user information. Let me know if there's anything else you'd like to add or adjust!