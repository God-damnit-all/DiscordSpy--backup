package commands

type Command interface {
  Invokes(); Description() string // Invokes is the command name, description is obvious
  Exec(ctx *Context) error // Execute command how you want to
}
