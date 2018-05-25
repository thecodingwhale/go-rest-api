## App
  - create a simple api for a job posting website.

### API

```
type User struct {
  Id
  Email
  Password
}
```
- [POST] | /users -> creating an account
  body -> {
    "email": "foo@email.com",
    "password": "password",
  }
```
type Job struct {
  Id
  UserId
  Name
  Description
}
```
* /jobs/${id} | protected
  * [GET] -

### Books
  https://www.golang-book.com/books/intro

go build && ./go-rest-api

for form validation
https://github.com/go-ozzo/ozzo-validation
