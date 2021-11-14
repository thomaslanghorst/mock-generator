# mock-generator

This is a small project I wrote to help create mock files using golangs [template package](https://pkg.go.dev/text/template). It will generate mock files that would be written using [the mock package of testify](https://github.com/stretchr/testify).

# Build the binary
To build the mock-generator you simply have to build the binary using
```
go build -o ./bin/mock-generator
```
If you want to make it gobally available, just export the path to your `~/.bashrc` file.

```
export PATH=$PATH:</path/to/your/project>/bin/mock-generator
```

# Running the mock-generator

After you have build the binary you can call and test it using the `./example/service.go` file.

```
./bin/mock-generator -i ./example/service.go
```
This will generate a file called `service_mock.go` next to the `service.go` file.

# Example
If you take this simple example 
```go
package service

type User struct {
}

type ServiceInterface interface {
	Login(userId string, password string) error
	ListUsers() ([]User, error)
	Logout(userId string)
}
```
and run the mock-generator. The output file will look like this:
```go
package service

import "github.com/stretchr/testify/mock"

type MockServiceInterface struct {
	mock.Mock
}

func (m *MockServiceInterface) Login(userId string, password string) error {
    args := m.Called(userId, password)
	
	return args.Error(0)
}

func (m *MockServiceInterface) ListUsers() ([]User, error) {
    args := m.Called()
	
	var v0 []User
	if args.Get(0) != nil {
		v0 = args.Get(0).([]User)
	}
	return v0, args.Error(1)
}

func (m *MockServiceInterface) Logout(userId string)  {
    m.Called(userId)
	
	
}
```
