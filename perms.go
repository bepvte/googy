package main

var perms map[string]map[string][]permrule

const (
Channel = iota + 1
Role
Server
)

type permrule struct {
	target, desc string
	targettype   int
	state        bool
}

