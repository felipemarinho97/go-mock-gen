# go-mock-gen

## Instalation

    $ git clone https://github.com/felipemarinho97/go-mock-gen.git

## Usage

Currently, is pretty rudimentar. Clone the code and edit the `GenerateMockCode` function, by instantiating the struct instance you want to mock in the `typeToMock` variable.

Also, edit the `Name` attribute in the `replacer` struct and write the intended interface name.

For example:

```go
	typeToMock := s3.Client{}
	replacer := Interface{
		Name: "S3Client",
	}
```

Will generate an output like that:

```go
// IS3Client generic client
type IS3Client interface {
        // AbortMultipartUpload godoc if available
    	AbortMultipartUpload(arg1 context.Context, arg2 *s3.AbortMultipartUploadInput, arg3 ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error)
    ........
```

Then, build and run the code.

```bash
    $ go run main.go > clients/s3_client.go
```
