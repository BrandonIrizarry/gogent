/*
Package logger implements a custom logger based on Go's built-in log package.

This package takes a different philosophy towards debug levels:
instead of using forcibly inclusive debug _levels_, it uses debug
_modes_ that can be individually toggled, which ultimately reflect the
programmer's intention to ignore or else exclusively focus on
particular areas of interest in the code execution.
*/
package logger
