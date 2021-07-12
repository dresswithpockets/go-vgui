package vgui

type ErrEmptyInput struct {}
type ErrEndOfInput struct {}

func (e *ErrEmptyInput) Error() string {
    return "The provided source input is empty."
}

func (e *ErrEndOfInput) Error() string {
    return "Cannot peek or advance past the end of the input"
}