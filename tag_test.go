package color

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTag(t *testing.T) {
	// force open color render for testing
	forceOpenColorRender()
	defer resetColorRender()

	is := assert.New(t)

	// sample 1
	r := String("<err>text</>")
	is.NotContains(r, "<")
	is.NotContains(r, ">")

	// disable color
	Enable = false
	r = Text("<err>text</>")
	is.Equal("text", r)
	Enable = true

	// sample 2
	s := "abc <err>err-text</> def <info>info text</>"
	r = ReplaceTag(s)
	is.NotContains(r, "<")
	is.NotContains(r, ">")

	// sample 3
	s = `abc <err>err-text</> 
def <info>info text
</>`
	r = ReplaceTag(s)
	is.NotContains(r, "<")
	is.NotContains(r, ">")

	// sample 4
	s = "abc <err>err-text</> def <err>err-text</> "
	r = ReplaceTag(s)
	is.NotContains(r, "<")
	is.NotContains(r, ">")

	// sample 5
	s = "abc <err>err-text</> def <d>"
	r = ReplaceTag(s)
	is.NotContains(r, "<err>")
	is.Contains(r, "<d>")

	// sample 6
	s = "custom tag: <fg=yellow;bg=black;op=underscore;>hello, welcome</>"
	r = ReplaceTag(s)
	is.NotContains(r, "<")
	is.NotContains(r, ">")

	s = Render()
	is.Equal("", s)
}

func TestParseCodeFromAttr(t *testing.T) {
	is := assert.New(t)

	s := ParseCodeFromAttr("=")
	is.Equal("", s)

	s = ParseCodeFromAttr("fg=lightRed;bg=lightRed;op=bold,blink")
	is.Equal("91;100;1;5", s)

	s = ParseCodeFromAttr("fg= lightRed;bg=lightRed;op=bold,")
	is.Equal("91;100;1", s)

	s = ParseCodeFromAttr("fg =lightRed;bg=lightRed;op=bold,blink")
	is.Equal("91;100;1;5", s)

	s = ParseCodeFromAttr("fg = lightRed;bg=lightRed;op=bold,blink")
	is.Equal("91;100;1;5", s)
}

func TestPrint(t *testing.T) {
	// force open color render for testing
	forceOpenColorRender()
	defer resetColorRender()
	is := assert.New(t)

	is.True(len(GetColorTags()) > 0)
	is.True(IsDefinedTag("info"))
	is.Equal("0;32", GetTagCode("info"))
	is.Equal("", GetTagCode("not-exist"))

	s := Sprint("<red>MSG</>")
	is.Equal("\x1b[0;31mMSG\x1b[0m", s)

	s = Sprint("<red>H</><green>I</>")
	is.Equal("\x1b[0;31mH\x1b[0m\x1b[0;32mI\x1b[0m", s)

	s = Sprintf("<red>%s</>", "MSG")
	is.Equal("\x1b[0;31mMSG\x1b[0m", s)

	// Print
	rewriteStdout()
	Print("<red>MSG</>")
	s = restoreStdout()
	is.Equal("\x1b[0;31mMSG\x1b[0m", s)

	// Printf
	rewriteStdout()
	Printf("<red>%s</>", "MSG")
	s = restoreStdout()
	is.Equal("\x1b[0;31mMSG\x1b[0m", s)

	// Println
	rewriteStdout()
	Println("<red>MSG</>")
	s = restoreStdout()
	is.Equal("\x1b[0;31mMSG\x1b[0m\n", s)

	rewriteStdout()
	Println("<red>hello</>", "world")
	s = restoreStdout()
	is.Equal("\x1b[0;31mhello\x1b[0m world\n", s)

	buf := new(bytes.Buffer)

	// Fprint
	Fprint(buf, "<red>MSG</>")
	is.Equal("\x1b[0;31mMSG\x1b[0m", buf.String())
	buf.Reset()

	// Fprintln
	_, err := Fprintln(buf, "<red>MSG</>")
	is.Equal("\x1b[0;31mMSG\x1b[0m\n", buf.String())
	is.NoError(err)
	buf.Reset()
	_, err = Fprintln(buf, "<red>hello</>", "world")
	is.Equal("\x1b[0;31mhello\x1b[0m world\n", buf.String())
	is.NoError(err)
	buf.Reset()

	// Fprintf
	_, err = Fprintf(buf, "<red>%s</>", "MSG")
	is.NoError(err)
	is.Equal("\x1b[0;31mMSG\x1b[0m", buf.String())
	buf.Reset()
}

func TestWrapTag(t *testing.T) {
	at := assert.New(t)
	at.Equal("<info>text</>", WrapTag("text", "info"))
	at.Equal("", WrapTag("", "info"))
	at.Equal("text", WrapTag("text", ""))
}

func TestApplyTag(t *testing.T) {
	forceOpenColorRender()
	defer resetColorRender()
	at := assert.New(t)
	at.Equal("\x1b[0;32mMSG\x1b[0m", ApplyTag("info", "MSG"))
}

func TestClearTag(t *testing.T) {
	is := assert.New(t)
	is.Equal("text", ClearTag("text"))
	is.Equal("text", ClearTag("<err>text</>"))
	is.Equal("abc error def info text", ClearTag("abc <err>error</> def <info>info text</>"))

	str := `abc <err>err-text</> 
def <info>info text
</>`
	ret := ClearTag(str)
	is.Contains(ret, "def info")
	is.NotContains(ret, "</>")

	str = "abc <err>text</> def<d>"
	ret = ClearTag(str)
	is.Equal("abc text def", ret)
	is.NotContains(ret, "<err>")
}

func TestTag_Print(t *testing.T) {
	forceOpenColorRender()
	defer resetColorRender()
	is := assert.New(t)

	s := Tag("info").Sprint("msg")
	is.Equal("\x1b[0;32mmsg\x1b[0m", s)

	s = Tag("info").Sprintf("m%s", "sg")
	is.Equal("\x1b[0;32mmsg\x1b[0m", s)

	info := Tag("info")

	// Tag.Print
	rewriteStdout()
	info.Print("msg")
	s = restoreStdout()
	if isLikeInCmd {
		is.Equal("msg", s)
	} else {
		is.Equal("\x1b[0;32mmsg\x1b[0m", s)
	}

	// Tag.Println
	rewriteStdout()
	info.Println("msg")
	s = restoreStdout()
	if isLikeInCmd {
		is.Equal("msg\n", s)
	} else {
		is.Equal("\x1b[0;32mmsg\x1b[0m\n", s)
	}

	// Tag.Printf
	rewriteStdout()
	info.Printf("m%s", "sg")
	s = restoreStdout()
	if isLikeInCmd {
		is.Equal("msg", s)
	} else {
		is.Equal("\x1b[0;32mmsg\x1b[0m", s)
	}

	mga := Tag("mga")

	// Tag.Print
	rewriteStdout()
	mga.Print("msg")
	s = restoreStdout()
	is.Equal("\x1b[0;35mmsg\x1b[0m", s)

	// Tag.Println
	rewriteStdout()
	mga.Println("msg")
	s = restoreStdout()
	is.Equal("\x1b[0;35mmsg\x1b[0m\n", s)

	// Tag.Printf
	rewriteStdout()
	mga.Printf("m%s", "sg")
	s = restoreStdout()
	is.Equal("\x1b[0;35mmsg\x1b[0m", s)
}
