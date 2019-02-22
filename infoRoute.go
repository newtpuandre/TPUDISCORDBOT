package main

//CommandRoute is used within AbsoluteRoute
type CommandRoute struct {
	Command string
}

//AbsoluteRoute is the struct returned when "/" is visited
type AbsoluteRoute struct {
	Message  string
	Commands []CommandRoute
}
