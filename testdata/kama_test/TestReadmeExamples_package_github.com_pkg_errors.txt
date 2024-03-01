
[
NAME: github.com/pkg/errors
CONSTANTS: []
VARIABLES: []
FUNCTIONS: [
	As(err error, target interface{}) bool 
	Cause(err error) error 
	Errorf(format string, args ...interface{}) error 
	Is(err error, target error) bool 
	New(message string) error 
	Unwrap(err error) error 
	WithMessage(err error, message string) error 
	WithMessagef(err error, format string, args ...interface{}) error 
	WithStack(err error) error 
	Wrap(err error, message string) error 
	Wrapf(err error, format string, args ...interface{}) error 
	]
TYPES: [
	Frame uintptr
		(Frame) Format(s fmt.State, verb rune)
		(Frame) MarshalText() ([]byte, error) 
	StackTrace []Frame
		(StackTrace) Format(s fmt.State, verb rune)]
]