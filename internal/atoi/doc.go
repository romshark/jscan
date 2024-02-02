// Package atoi provides integer parsing functions optimized specifically for
// jscan since the parser validates the input before the integer parser is invoked.
//
// WARNING: These functions are not a replacement for strconv.ParseInt or strconv.Atoi
// because they don't validate the input and assumes valid input instead.
// The jscan tokenizer is guaranteed to provide only valid values which
// are only not guaranteed to not overflow.
package atoi
