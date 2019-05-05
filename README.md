# Quick intro

A tool for managing chunks of static data in Go programs.

## install & usage

Install with
```
> go get github.com/nyarly/inlinefiles
```
In your code:
```go
//go:generate inlinefiles --vfs=Json json_blobs lib/jsonblobs.go

func getJSON(path string) string {
  buf := &bytes.Buffer{}
  json, _ := Json.Open(path)
  io.Copy(buf, json)
  return buf.String()
}
```
(this code is purely for example purposes.)

Then
```
> go generate
```
which will create a `lib/jsonblobs.go` file
with the `Json` constant holding a `mapfs` flavored `vfs.FileSystem`.


# inlinefiles

Say you've got a simple web-service to write.
It's going to be a microservice,
and in order to
play nice in its environment,
you want it to serve a static file
with JSON that describes
where it's monitoring page is
and who to get in touch with if it misbehaves.
You could
include a literal JSON string,
but it becomes harder to work with
than it would be in its own file.

In its own file,
your editor can highlight the syntax correctly,
and check that the picky JSON syntax is correct.
The problem is more significant
if you want to use complicated templates,
where the concerns of whitespace
for the template output
and for the code itself
are in tension.

When there's a tree of template files,
where the templates include one another
and so they need a file heirarchy provided to them,
it get's a little more annoying to manage.
It's not impossible though -
you gin up a bunch of path strings and plug them in.

What we want
it to be able to pull arbitrary files
into Go code as string literals.
We'd like to be able to do this
as part of our build process,
because we want updates to the source files
to update the Go strings as automatically as possible.

Go has `go:generate`,
which is very handy for this.
If we include a magic comment in our code,
then running `go generate` on the command line
will run commands to produce code.
There's support in the stdlib
for generating `Stringer` interfaces,
for instance.

It would also be handy
if we could continue to treat the files
_as files_ instead of strings.
In part,
because this means that it's easier to work
with template inclusion.
Also because it means that if we later want
to override the hardcoded files,
we're still using an `os.File` or `io.Buffer`
instead of mixing those types with `string`.

Go has the `vfs` package.
It's part of `go doc`
but it serves our purposes perfectly.
You can define a `mapfs`
which is a literal `map[string]string`,
but allows you to `Open(path)` and get file structs.

So there's a lot of facility to do what we want in go already.
All that's missing is a tool to take
a real directory structure,
and generate the code for a literal

Here's `inlinefiles` to weld all this together.
The design intention was
to adhere to the Unix principle
of small tools doing one thing well.
`inlinefiles` takes two arguments:
a path to a directory of source files,
and a path to an output file,
and generates Go code to build an equivalent `mapfs`.
There's a few flags to handle things like
using a specific package name (instead of guessing)
and to restrict the files included based on a glob pattern.
It's suiteable for using as a `//go:generate` comment
(since that was the motivating case.)

## templatestore package

In addition, `inlinefiles` includes a package for
smoothing out using `vfs.FileSystem`s with
Go templates.
It's very lightweight,
but if you `import github.com/nyarly/inlinefiles/templatestore`
you get `LoadText` and `LoadTextOnto` which'll let you
construct a `template.Template` with subtemplates suitable
for all your most demanding templating applications.
