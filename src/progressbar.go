package src

import (
	"errors"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/colorstring"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a thread-safe, simple
// progress bar
type ProgressBar struct {
	state  state
	config configProgressbar
	lock   sync.Mutex
}

// State is the basic properties of the bar
type State struct {
	CurrentPercent float64
	CurrentBytes   float64
	SecondsSince   float64
	SecondsLeft    float64
	KBsPerSecond   float64
}

type state struct {
	currentNum        int64
	currentPercent    int
	lastPercent       int
	currentSaucerSize int
	isAltSaucerHead   bool

	lastShown time.Time
	startTime time.Time

	counterTime         time.Time
	counterNumSinceLast int64
	counterLastTenRates []float64

	maxLineWidth int
	currentBytes float64
	finished     bool

	rendered string
}

type configProgressbar struct {
	max                  int64 // max number of the counter
	maxHumanized         string
	maxHumanizedSuffix   string
	width                int
	writer               io.Writer
	theme                Theme
	renderWithBlankState bool
	description          string
	iterationString      string
	ignoreLength         bool // ignoreLength if max bytes not known

	// whether the output is expected to contain color codes
	colorCodes bool

	// show rate of change in kB/sec or MB/sec
	showBytes bool
	// show the iterations per second
	showIterationsPerSecond bool
	showIterationsCount     bool

	// whether the progress bar should attempt to predict the finishing
	// time of the progress based on the start time and the average
	// number of seconds between  increments.
	predictTime bool

	// minimum time to wait in between updates
	throttleDuration time.Duration

	// clear bar once finished
	clearOnFinish bool

	// spinnerType should be a number between 0-75
	spinnerType int

	// fullWidth specifies whether to measure and set the bar to a specific width
	fullWidth bool

	// invisible doesn't render the bar at all, useful for debugging
	invisible bool

	onCompletion func()

	// whether the render function should make use of ANSI codes to reduce console I/O
	useANSICodes bool
}

// Theme defines the elements of the bar
type Theme struct {
	Saucer        string
	AltSaucerHead string
	SaucerHead    string
	SaucerPadding string
	BarStart      string
	BarEnd        string
}

// Option is the type all options need to adhere to
type Option func(p *ProgressBar)

var defaultTheme = Theme{Saucer: "█", SaucerPadding: " ", BarStart: "|", BarEnd: "|"}

// NewOptions constructs a new instance of ProgressBar, with any options you specify
func NewOptions(max int, options ...Option) *ProgressBar {
	return NewOptions64(int64(max), options...)
}

// NewOptions64 constructs a new instance of ProgressBar, with any options you specify
func NewOptions64(max int64, options ...Option) *ProgressBar {
	b := ProgressBar{
		state: getBasicState(),
		config: configProgressbar{
			writer:           os.Stdout,
			theme:            defaultTheme,
			iterationString:  "it",
			width:            40,
			max:              max,
			throttleDuration: 0 * time.Nanosecond,
			predictTime:      true,
			spinnerType:      9,
			invisible:        false,
		},
	}

	for _, o := range options {
		o(&b)
	}

	if b.config.spinnerType < 0 || b.config.spinnerType > 75 {
		panic("invalid spinner type, must be between 0 and 75")
	}

	// ignoreLength if max bytes not known
	if b.config.max == -1 {
		b.config.ignoreLength = true
		b.config.max = int64(b.config.width)
		b.config.predictTime = false
	}

	b.config.maxHumanized, b.config.maxHumanizedSuffix = humanizeBytes(float64(b.config.max))

	if b.config.renderWithBlankState {

		if err := b.RenderBlank(); err != nil {
			log.Println(err)
		}
	}

	return &b
}

func getBasicState() state {
	now := time.Now()
	return state{
		startTime:   now,
		lastShown:   now,
		counterTime: now,
	}
}

// New returns a new ProgressBar
// with the specified maximum
func New(max int) *ProgressBar {
	return NewOptions(max)
}

// String returns the current rendered version of the progress bar.
// It will never return an empty string while the progress bar is running.
func (p *ProgressBar) String() string {
	return p.state.rendered
}

// RenderBlank renders the current bar state, you can use this to render a 0% state
func (p *ProgressBar) RenderBlank() error {
	if p.config.invisible {
		return nil
	}
	return p.render()
}

// Reset will reset the clock that is used
// to calculate current time and the time left.
func (p *ProgressBar) Reset() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.state = getBasicState()
}

// Finish will fill the bar to full
func (p *ProgressBar) Finish() error {
	p.lock.Lock()
	p.state.currentNum = p.config.max
	p.lock.Unlock()
	return p.Add(0)
}

// Add will add the specified amount to the progressbar
func (p *ProgressBar) Add(num int) error {
	return p.Add64(int64(num))
}

// Set wil set the bar to a current number
func (p *ProgressBar) Set(num int) error {
	return p.Set64(int64(num))
}

// Set64 wil set the bar to a current number
func (p *ProgressBar) Set64(num int64) error {
	p.lock.Lock()
	toAdd := num - int64(p.state.currentBytes)
	p.lock.Unlock()
	return p.Add64(toAdd)
}

// Add64 will add the specified amount to the progressbar
func (p *ProgressBar) Add64(num int64) error {
	if p.config.invisible {
		return nil
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.config.max == 0 {
		return errors.New("max must be greater than 0")
	}

	if p.state.currentNum < p.config.max {
		if p.config.ignoreLength {
			p.state.currentNum = (p.state.currentNum + num) % p.config.max
		} else {
			p.state.currentNum += num
		}
	}

	p.state.currentBytes += float64(num)

	// reset the countdown timer every second to take rolling average
	p.state.counterNumSinceLast += num
	if time.Since(p.state.counterTime).Seconds() > 0.5 {
		p.state.counterLastTenRates = append(p.state.counterLastTenRates, float64(p.state.counterNumSinceLast)/time.Since(p.state.counterTime).Seconds())
		if len(p.state.counterLastTenRates) > 10 {
			p.state.counterLastTenRates = p.state.counterLastTenRates[1:]
		}
		p.state.counterTime = time.Now()
		p.state.counterNumSinceLast = 0
	}

	percent := float64(p.state.currentNum) / float64(p.config.max)
	p.state.currentSaucerSize = int(percent * float64(p.config.width))
	p.state.currentPercent = int(percent * 100)
	updateBar := p.state.currentPercent != p.state.lastPercent && p.state.currentPercent > 0

	p.state.lastPercent = p.state.currentPercent
	if p.state.currentNum > p.config.max {
		return errors.New("current number exceeds max")
	}

	// always update if show bytes/second or its/second
	if updateBar || p.config.showIterationsPerSecond || p.config.showIterationsCount {
		return p.render()
	}

	return nil
}

// Clear erases the progress bar from the current line
func (p *ProgressBar) Clear() error {
	return clearProgressBar(p.config, p.state)
}

// Describe will change the description shown before the progress, which
// can be changed on the fly (as for a slow running process).
func (p *ProgressBar) Describe(description string) {
	p.config.description = description
	if err := p.RenderBlank(); err != nil {
		log.Println(err)
	}
}

// GetMax returns the max of a bar
func (p *ProgressBar) GetMax() int {
	return int(p.config.max)
}

// GetMax64 returns the current max
func (p *ProgressBar) GetMax64() int64 {
	return p.config.max
}

// ChangeMax takes in int
// and changes the max value
// of the progress bar
func (p *ProgressBar) ChangeMax(newMax int) {
	p.ChangeMax64(int64(newMax))
}

// ChangeMax64 is basically
// the same as ChangeMax,
// but takes in int64
// to avoid casting
func (p *ProgressBar) ChangeMax64(newMax int64) {
	p.config.max = newMax

	if p.config.showBytes {
		p.config.maxHumanized, p.config.maxHumanizedSuffix = humanizeBytes(float64(p.config.max))
	}

	if err := p.Add(0); err != nil { // re-render
		log.Println(err)
	}
}

// IsFinished returns true if progress bar is completed
func (p *ProgressBar) IsFinished() bool {
	return p.state.finished
}

// render renders the progress bar, updating the maximum
// rendered line width. this function is not thread-safe,
// so it must be called with an acquired lock.
func (p *ProgressBar) render() error {
	// make sure that the rendering is not happening too quickly
	// but always show if the currentNum reaches the max
	if time.Since(p.state.lastShown).Nanoseconds() < p.config.throttleDuration.Nanoseconds() &&
		p.state.currentNum < p.config.max {
		return nil
	}

	if !p.config.useANSICodes {
		// first, clear the existing progress bar
		err := clearProgressBar(p.config, p.state)
		if err != nil {
			return err
		}
	}

	// check if the progress bar is finished
	if !p.state.finished && p.state.currentNum >= p.config.max {
		p.state.finished = true
		if !p.config.clearOnFinish {
			if _, err := renderProgressBar(p.config, &p.state); err != nil {
				log.Println(err)
			}
		}

		if p.config.onCompletion != nil {
			p.config.onCompletion()
		}
	}
	if p.state.finished {
		// when using ANSI codes we don't pre-clean the current line
		if p.config.useANSICodes {
			err := clearProgressBar(p.config, p.state)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// then, re-render the current progress bar
	w, err := renderProgressBar(p.config, &p.state)
	if err != nil {
		return err
	}

	if w > p.state.maxLineWidth {
		p.state.maxLineWidth = w
	}

	p.state.lastShown = time.Now()

	return nil
}

// State returns the current state
func (p *ProgressBar) State() State {
	p.lock.Lock()
	defer p.lock.Unlock()
	s := State{}
	s.CurrentPercent = float64(p.state.currentNum) / float64(p.config.max)
	s.CurrentBytes = p.state.currentBytes
	s.SecondsSince = time.Since(p.state.startTime).Seconds()
	if p.state.currentNum > 0 {
		s.SecondsLeft = s.SecondsSince / float64(p.state.currentNum) * (float64(p.config.max) - float64(p.state.currentNum))
	}
	s.KBsPerSecond = float64(p.state.currentBytes) / 1024.0 / s.SecondsSince
	return s
}

// regex matching ansi escape codes
var ansiRegex = regexp.MustCompile(`\x1b\[[\d;]*[a-zA-Z]`)

func getStringWidth(c configProgressbar, str string) int {
	if c.colorCodes {
		// convert any color codes in the progress bar into the respective ANSI codes
		str = colorstring.Color(str)
	}

	// the width of the string, if printed to the console
	// does not include the carriage return character
	cleanString := strings.ReplaceAll(str, "\r", "")

	if c.colorCodes {
		// the ANSI codes for the colors do not take up space in the console output,
		// so they do not count towards the output string width
		cleanString = ansiRegex.ReplaceAllString(cleanString, "")
	}

	// get the amount of runes in the string instead of the
	// character count of the string, as some runes span multiple characters.
	// see https://stackoverflow.com/a/12668840/2733724
	stringWidth := runewidth.StringWidth(cleanString)
	return stringWidth
}

func renderProgressBar(c configProgressbar, s *state) (int, error) {
	leftBrac := ""
	rightBrac := ""
	saucer := ""
	saucerHead := ""
	bytesString := ""
	str := ""

	averageRate := average(s.counterLastTenRates)
	if len(s.counterLastTenRates) == 0 || s.finished {
		// if no average samples, or if finished,
		// then average rate should be the total rate
		averageRate = s.currentBytes / time.Since(s.startTime).Seconds()
	}

	// show iteration count in "current/total" iterations format
	if c.showIterationsCount {
		bytesString += "("
		if !c.ignoreLength {
			if c.showBytes {
				currentHumanize, currentSuffix := humanizeBytes(s.currentBytes)
				if currentSuffix == c.maxHumanizedSuffix {
					bytesString += fmt.Sprintf("%s/%s%s", currentHumanize, c.maxHumanized, c.maxHumanizedSuffix)
				} else {
					bytesString += fmt.Sprintf("%s%s/%s%s", currentHumanize, currentSuffix, c.maxHumanized, c.maxHumanizedSuffix)
				}
			} else {
				bytesString += fmt.Sprintf("%.0f/%d", s.currentBytes, c.max)
			}
		} else {
			if c.showBytes {
				currentHumanize, currentSuffix := humanizeBytes(s.currentBytes)
				bytesString += fmt.Sprintf("%s%s", currentHumanize, currentSuffix)
			} else {
				bytesString += fmt.Sprintf("%.0f/%s", s.currentBytes, "-")
			}
		}
	}

	// show rolling average rate in kB/sec or MB/sec
	if c.showBytes {
		if bytesString == "" {
			bytesString += "("
		} else {
			bytesString += ", "
		}
		kbPerSecond := averageRate / 1024.0
		if kbPerSecond > 1024.0 {
			bytesString += fmt.Sprintf("%0.3f MB/s", kbPerSecond/1024.0)
		} else if kbPerSecond > 0 {
			bytesString += fmt.Sprintf("%0.3f kB/s", kbPerSecond)
		}
	}

	// show iterations rate
	if c.showIterationsPerSecond {
		if bytesString == "" {
			bytesString += "("
		} else {
			bytesString += ", "
		}
		if averageRate > 1 {
			bytesString += fmt.Sprintf("%0.0f %s/s", averageRate, c.iterationString)
		} else {
			bytesString += fmt.Sprintf("%0.0f %s/min", 60*averageRate, c.iterationString)
		}
	}
	if bytesString != "" {
		bytesString += ")"
	}

	// show time prediction in "current/total" seconds format
	if c.predictTime {
		leftBrac = (time.Duration(time.Since(s.startTime).Seconds()) * time.Second).String()
		rightBracNum := time.Duration((1/averageRate)*(float64(c.max)-float64(s.currentNum))) * time.Second
		if rightBracNum.Seconds() < 0 {
			rightBracNum = 0 * time.Second
		}
		rightBrac = rightBracNum.String()
	}

	if c.fullWidth && !c.ignoreLength {
		width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width, _, err = terminal.GetSize(int(os.Stderr.Fd()))
			if err != nil {
				width = 80
			}
		}

		c.width = width - getStringWidth(c, c.description) - 14 - len(bytesString) - len(leftBrac) - len(rightBrac)
		s.currentSaucerSize = int(float64(s.currentPercent) / 100.0 * float64(c.width))
	}
	if s.currentSaucerSize > 0 {
		if c.ignoreLength {
			saucer = strings.Repeat(c.theme.SaucerPadding, s.currentSaucerSize-1)
		} else {
			saucer = strings.Repeat(c.theme.Saucer, s.currentSaucerSize-1)
		}

		// Check if an alternate saucer head is set for animation
		if c.theme.AltSaucerHead != "" && s.isAltSaucerHead {
			saucerHead = c.theme.AltSaucerHead
			s.isAltSaucerHead = false
		} else if c.theme.SaucerHead == "" || s.currentSaucerSize == c.width {
			// use the saucer for the saucer head if it hasn't been set
			// to preserve backwards compatibility
			saucerHead = c.theme.Saucer
		} else {
			saucerHead = c.theme.SaucerHead
			s.isAltSaucerHead = true
		}
		saucer += saucerHead
	}

	/*
		Progress Bar format
		Description % |------        |  (kb/s) (iteration count) (iteration rate) (predict time)
	*/
	repeatAmount := c.width - s.currentSaucerSize
	if repeatAmount < 0 {
		repeatAmount = 0
	}
	if c.ignoreLength {
		str = fmt.Sprintf("\r%s %s %s ",
			spinners[c.spinnerType][int(math.Round(math.Mod(float64(time.Since(s.startTime).Milliseconds()/100), float64(len(spinners[c.spinnerType])))))],
			c.description,
			bytesString,
		)
	} else if leftBrac == "" {
		str = fmt.Sprintf("\r%s%4d%% %s%s%s%s %s ",
			c.description,
			s.currentPercent,
			c.theme.BarStart,
			saucer,
			strings.Repeat(c.theme.SaucerPadding, repeatAmount),
			c.theme.BarEnd,
			bytesString,
		)
	} else {
		if s.currentPercent == 100 {
			str = fmt.Sprintf("\r%s%4d%% %s%s%s%s %s",
				c.description,
				s.currentPercent,
				c.theme.BarStart,
				saucer,
				strings.Repeat(c.theme.SaucerPadding, repeatAmount),
				c.theme.BarEnd,
				bytesString,
			)
		} else {
			str = fmt.Sprintf("\r%s%s%s%s%s %s%4d%% ",
				c.description,
				c.theme.BarStart,
				saucer,
				strings.Repeat(c.theme.SaucerPadding, repeatAmount),
				c.theme.BarEnd,
				bytesString,
				s.currentPercent,
			)
		}
	}

	if c.colorCodes {
		// convert any color codes in the progress bar into the respective ANSI codes
		str = colorstring.Color(str)
	}

	s.rendered = str

	return getStringWidth(c, str), writeString(c, str)
}

func clearProgressBar(c configProgressbar, s state) error {
	if c.useANSICodes {
		// write the "clear current line" ANSI escape sequence
		return writeString(c, "\033[2K\r")
	}
	// fill the empty content
	// to overwrite the progress bar and jump
	// back to the beginning of the line
	str := fmt.Sprintf("\r%s\r", strings.Repeat(" ", s.maxLineWidth))
	return writeString(c, str)
	// the following does not show correctly if the previous line is longer than subsequent line
	// return writeString(c, "\r")
}

func writeString(c configProgressbar, str string) error {
	if _, err := io.WriteString(c.writer, str); err != nil {
		return err
	}

	if f, ok := c.writer.(*os.File); ok {
		// ignore any errors in Sync(), as stdout
		// can't be synced on some operating systems
		// like Debian 9 (Stretch)
		if err := f.Sync(); err != nil {
		}
	}

	return nil
}

// Reader is the progressbar io.Reader struct
type Reader struct {
	io.Reader
	bar *ProgressBar
}

// Read will read the data and add the number of bytes to the progressbar
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err = r.bar.Add(n); err != nil {
		log.Println(err)
	}
	return
}

// Close the reader when it implements io.Closer
func (r *Reader) Close() (err error) {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	if err = r.bar.Finish(); err != nil {
		log.Println(err)
	}
	return
}

// Write implement io.Writer
func (p *ProgressBar) Write(b []byte) (n int, err error) {
	n = len(b)
	if err = p.Add(n); err != nil {
		log.Println(err)
	}
	return
}

// Read implement io.Reader
func (p *ProgressBar) Read(b []byte) (n int, err error) {
	n = len(b)
	if err = p.Add(n); err != nil {
		log.Println(err)
	}
	return
}

func (p *ProgressBar) Close() (err error) {
	if err = p.Finish(); err != nil {
		log.Println(err)
	}
	return
}

func average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

func humanizeBytes(s float64) (string, string) {
	sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
	base := 1024.0
	if s < 10 {
		return fmt.Sprintf("%2.0f", s), "B"
	}
	e := math.Floor(logn(s, base))
	suffix := sizes[int(e)]
	val := math.Floor(s/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

var spinners = map[int][]string{
	0:  {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	1:  {"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▁"},
	2:  {"▖", "▘", "▝", "▗"},
	3:  {"┤", "┘", "┴", "└", "├", "┌", "┬", "┐"},
	4:  {"◢", "◣", "◤", "◥"},
	5:  {"◰", "◳", "◲", "◱"},
	6:  {"◴", "◷", "◶", "◵"},
	7:  {"◐", "◓", "◑", "◒"},
	8:  {".", "o", "O", "@", "*"},
	9:  {"|", "/", "-", "\\"},
	10: {"◡◡", "⊙⊙", "◠◠"},
	11: {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
	12: {">))'>", " >))'>", "  >))'>", "   >))'>", "    >))'>", "   <'((<", "  <'((<", " <'((<"},
	13: {"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
	14: {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	15: {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
	16: {"▉", "▊", "▋", "▌", "▍", "▎", "▏", "▎", "▍", "▌", "▋", "▊", "▉"},
	17: {"■", "□", "▪", "▫"},
	18: {"←", "↑", "→", "↓"},
	19: {"╫", "╪"},
	20: {"⇐", "⇖", "⇑", "⇗", "⇒", "⇘", "⇓", "⇙"},
	21: {"⠁", "⠁", "⠉", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠤", "⠄", "⠄", "⠤", "⠠", "⠠", "⠤", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋", "⠉", "⠈", "⠈"},
	22: {"⠈", "⠉", "⠋", "⠓", "⠒", "⠐", "⠐", "⠒", "⠖", "⠦", "⠤", "⠠", "⠠", "⠤", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋", "⠉", "⠈"},
	23: {"⠁", "⠉", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠤", "⠄", "⠄", "⠤", "⠴", "⠲", "⠒", "⠂", "⠂", "⠒", "⠚", "⠙", "⠉", "⠁"},
	24: {"⠋", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋"},
	25: {"ｦ", "ｧ", "ｨ", "ｩ", "ｪ", "ｫ", "ｬ", "ｭ", "ｮ", "ｯ", "ｱ", "ｲ", "ｳ", "ｴ", "ｵ", "ｶ", "ｷ", "ｸ", "ｹ", "ｺ", "ｻ", "ｼ", "ｽ", "ｾ", "ｿ", "ﾀ", "ﾁ", "ﾂ", "ﾃ", "ﾄ", "ﾅ", "ﾆ", "ﾇ", "ﾈ", "ﾉ", "ﾊ", "ﾋ", "ﾌ", "ﾍ", "ﾎ", "ﾏ", "ﾐ", "ﾑ", "ﾒ", "ﾓ", "ﾔ", "ﾕ", "ﾖ", "ﾗ", "ﾘ", "ﾙ", "ﾚ", "ﾛ", "ﾜ", "ﾝ"},
	26: {".", "..", "..."},
	27: {"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▉", "▊", "▋", "▌", "▍", "▎", "▏", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁"},
	28: {".", "o", "O", "°", "O", "o", "."},
	29: {"+", "x"},
	30: {"v", "<", "^", ">"},
	31: {">>--->", " >>--->", "  >>--->", "   >>--->", "    >>--->", "    <---<<", "   <---<<", "  <---<<", " <---<<", "<---<<"},
	32: {"|", "||", "|||", "||||", "|||||", "|||||||", "||||||||", "|||||||", "||||||", "|||||", "||||", "|||", "||", "|"},
	33: {"[          ]", "[=         ]", "[==        ]", "[===       ]", "[====      ]", "[=====     ]", "[======    ]", "[=======   ]", "[========  ]", "[========= ]", "[==========]"},
	34: {"(*---------)", "(-*--------)", "(--*-------)", "(---*------)", "(----*-----)", "(-----*----)", "(------*---)", "(-------*--)", "(--------*-)", "(---------*)"},
	35: {"█▒▒▒▒▒▒▒▒▒", "███▒▒▒▒▒▒▒", "█████▒▒▒▒▒", "███████▒▒▒", "██████████"},
	36: {"[                    ]", "[=>                  ]", "[===>                ]", "[=====>              ]", "[======>             ]", "[========>           ]", "[==========>         ]", "[============>       ]", "[==============>     ]", "[================>   ]", "[==================> ]", "[===================>]"},
	37: {"ဝ", "၀"},
	38: {"▌", "▀", "▐▄"},
	39: {"🌍", "🌎", "🌏"},
	40: {"◜", "◝", "◞", "◟"},
	41: {"⬒", "⬔", "⬓", "⬕"},
	42: {"⬖", "⬘", "⬗", "⬙"},
	43: {"[>>>          >]", "[]>>>>        []", "[]  >>>>      []", "[]    >>>>    []", "[]      >>>>  []", "[]        >>>>[]", "[>>          >>]"},
	44: {"♠", "♣", "♥", "♦"},
	45: {"➞", "➟", "➠", "➡", "➠", "➟"},
	46: {"  |  ", ` \   `, "_    ", ` \   `, "  |  ", "   / ", "    _", "   / "},
	47: {"  . . . .", ".   . . .", ". .   . .", ". . .   .", ". . . .  ", ". . . . ."},
	48: {" |     ", "  /    ", "   _   ", `    \  `, "     | ", `    \  `, "   _   ", "  /    "},
	49: {"⎺", "⎻", "⎼", "⎽", "⎼", "⎻"},
	50: {"▹▹▹▹▹", "▸▹▹▹▹", "▹▸▹▹▹", "▹▹▸▹▹", "▹▹▹▸▹", "▹▹▹▹▸"},
	51: {"[    ]", "[   =]", "[  ==]", "[ ===]", "[====]", "[=== ]", "[==  ]", "[=   ]"},
	52: {"( ●    )", "(  ●   )", "(   ●  )", "(    ● )", "(     ●)", "(    ● )", "(   ●  )", "(  ●   )", "( ●    )"},
	53: {"✶", "✸", "✹", "✺", "✹", "✷"},
	54: {"▐|\\____________▌", "▐_|\\___________▌", "▐__|\\__________▌", "▐___|\\_________▌", "▐____|\\________▌", "▐_____|\\_______▌", "▐______|\\______▌", "▐_______|\\_____▌", "▐________|\\____▌", "▐_________|\\___▌", "▐__________|\\__▌", "▐___________|\\_▌", "▐____________|\\▌", "▐____________/|▌", "▐___________/|_▌", "▐__________/|__▌", "▐_________/|___▌", "▐________/|____▌", "▐_______/|_____▌", "▐______/|______▌", "▐_____/|_______▌", "▐____/|________▌", "▐___/|_________▌", "▐__/|__________▌", "▐_/|___________▌", "▐/|____________▌"},
	55: {"▐⠂       ▌", "▐⠈       ▌", "▐ ⠂      ▌", "▐ ⠠      ▌", "▐  ⡀     ▌", "▐  ⠠     ▌", "▐   ⠂    ▌", "▐   ⠈    ▌", "▐    ⠂   ▌", "▐    ⠠   ▌", "▐     ⡀  ▌", "▐     ⠠  ▌", "▐      ⠂ ▌", "▐      ⠈ ▌", "▐       ⠂▌", "▐       ⠠▌", "▐       ⡀▌", "▐      ⠠ ▌", "▐      ⠂ ▌", "▐     ⠈  ▌", "▐     ⠂  ▌", "▐    ⠠   ▌", "▐    ⡀   ▌", "▐   ⠠    ▌", "▐   ⠂    ▌", "▐  ⠈     ▌", "▐  ⠂     ▌", "▐ ⠠      ▌", "▐ ⡀      ▌", "▐⠠       ▌"},
	56: {"¿", "?"},
	57: {"⢹", "⢺", "⢼", "⣸", "⣇", "⡧", "⡗", "⡏"},
	58: {"⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"},
	59: {".  ", ".. ", "...", " ..", "  .", "   "},
	60: {".", "o", "O", "°", "O", "o", "."},
	61: {"▓", "▒", "░"},
	62: {"▌", "▀", "▐", "▄"},
	63: {"⊶", "⊷"},
	64: {"▪", "▫"},
	65: {"□", "■"},
	66: {"▮", "▯"},
	67: {"-", "=", "≡"},
	68: {"d", "q", "p", "b"},
	69: {"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"},
	70: {"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "},
	71: {"☗", "☖"},
	72: {"⧇", "⧆"},
	73: {"◉", "◎"},
	74: {"㊂", "㊀", "㊁"},
	75: {"⦾", "⦿"},
}
