package main

//CommandRoute is used within AbsoluteRoute
type CommandRoute struct {
	Command string `json:""`
}

//AbsoluteRoute is the struct returned when "/" is visited
type AbsoluteRoute struct {
	Message  string         `json:"Message"`
	Commands []CommandRoute `json:"Available Commands"`
}
