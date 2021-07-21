package epilog

// Appender is an interface
type Appender interface {
	Append() (err error)
}

// Item is an interface.
// any type implements this interface
// can be used as parameter of Persist
type Item interface {
	Serialize() (res string)
}
