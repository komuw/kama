
[
NAME: net/http.Request
KIND: struct
SIGNATURE: [*http.Request http.Request]
FIELDS: [
	Method string 
	URL *url.URL 
	Proto string 
	ProtoMajor int 
	ProtoMinor int 
	Header http.Header 
	Body io.ReadCloser 
	GetBody func() (io.ReadCloser, error) 
	ContentLength int64 
	TransferEncoding []string 
	Close bool 
	Host string 
	Form url.Values 
	PostForm url.Values 
	MultipartForm *multipart.Form 
	Trailer http.Header 
	RemoteAddr string 
	RequestURI string 
	TLS *tls.ConnectionState 
	Cancel <-chan struct {} 
	Response *http.Response 
	]
METHODS: [
	AddCookie func(*http.Request, *http.Cookie) 
	BasicAuth func(*http.Request) (string, string, bool) 
	Clone func(*http.Request, context.Context) *http.Request 
	Context func(*http.Request) context.Context 
	Cookie func(*http.Request, string) (*http.Cookie, error) 
	Cookies func(*http.Request) []*http.Cookie 
	FormFile func(*http.Request, string) (multipart.File, *multipart.FileHeader, error) 
	FormValue func(*http.Request, string) string 
	MultipartReader func(*http.Request) (*multipart.Reader, error) 
	ParseForm func(*http.Request) error 
	ParseMultipartForm func(*http.Request, int64) error 
	PathValue func(*http.Request, string) string 
	PostFormValue func(*http.Request, string) string 
	ProtoAtLeast func(*http.Request, int, int) bool 
	Referer func(*http.Request) string 
	SetBasicAuth func(*http.Request, string, string) 
	SetPathValue func(*http.Request, string, string) 
	UserAgent func(*http.Request) string 
	WithContext func(*http.Request, context.Context) *http.Request 
	Write func(*http.Request, io.Writer) error 
	WriteProxy func(*http.Request, io.Writer) error 
	]
STACK_TRACE: [
[10mLEGEND:
 compiler: blue
 thirdParty: yellow
 yours: red

[0m[131m	github.com/komuw/kama/kama.go:106 github.com/komuw/kama.Dir
[0m[131m	github.com/komuw/kama/kama_test.go:246 github.com/komuw/kama.TestReadmeExamples.func1
[0m[131m	testing/testing.go:1689 testing.tRunner
[0m[131m	runtime/asm_amd64.s:1695 runtime.goexit
[0m]
SNIPPET: &Request{
  Method: "GET",
  URL: &URL{
    Scheme: "https",
    Host: "example.com",
  },
  Proto: "HTTP/1.1",
  ProtoMajor: int(1),
  ProtoMinor: int(1),
  Header: http.Header{
   "Content-Type": []string{
   "application/octet-stream",
      }, 
   "Cookie": []string{
   "hello=world",
      }, 
    },
  Body: io.ReadCloser nil,
  GetBody: func() (io.ReadCloser, error),
  ContentLength: int64(0),
  TransferEncoding: []string{(nil)},
  Close: false,
  Host: "example.com",
  Form: url.Values{(nil)},
  PostForm: url.Values{(nil)},
  MultipartForm: *multipart.Form(nil),
  Trailer: http.Header{(nil)},
  RemoteAddr: "",
  RequestURI: "",
  TLS: *tls.ConnectionState(nil),
  Cancel: <-chan struct {} (len=0, cap=0),
  Response: *http.Response(nil),
}
]
