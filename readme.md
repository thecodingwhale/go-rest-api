## App
  - create a simple api for a job posting with jwt authentication

### Run
```
go build && ./go-rest-api
```

type User struct {
  Id
  Email
  Name
  Password
}

### API
#### POST - /users
```
body
{
  "email": "john.doe@email.com",
  "name": "John Doe",
  "password": "12345678"
}
response
{}
```

#### POST - /authenticate
```
body
{
  "email": "john.doe@email.com",
  "password": "12345678"
}
response
{
  "token": "token"
}
```

#### POST - /jobs
```
headers:
  Key: Authorization
  Value: Bearer token
body
{
  "post": "Senior Software Engineer",
  "company": "A New World",
  "location": "Remote"
}
response
{}
```

#### GET - /test
```
headers:
  Key: Authorization
  Value: Bearer token
body
{
  "email": "john.doe@email.com",
  "password": "12345678"
}
response
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@email.com"
}
```

### Helpful Links
  https://www.golang-book.com/books/intro
  https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/05.2.html
  https://aligoren.com/golang-get-structtag-values/
  https://stackoverflow.com/questions/14514312/obtaining-the-name-of-a-known-struct-field
  https://stackoverflow.com/questions/18930910/access-struct-property-by-name
  https://golang.org/src/database/sql/example_test.go
  http://go-database-sql.org/retrieving.html
  https://blog.questionable.services/article/http-handler-error-handling-revisited/
  https://stackoverflow.com/questions/6012692/os-error-string-value-golang
  https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee
  https://gowebexamples.com/password-hashing/
  https://github.com/go-ozzo/ozzo-validation

