package msg

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"strings"
	"time"

	"runtime/debug"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gookit/event"
)

const (
	HEADR = `欢迎弹幕指导☝️ 🤓`
)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

func Start() {
	p := tea.NewProgram(initialModel())
	listenMsg(p)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type DanMuMsg struct {
	Tag     string
	From    string
	Content string
	Type    int
}

// 简单的情况就不要用opts模式了
func NewDanMuMsg(tag, from, content string, dType int) *DanMuMsg {
	return &DanMuMsg{
		Tag:     tag,
		From:    from,
		Content: content,
		Type:    dType,
	}
}

type resultMsg struct {
	recvTime string
	content  string
}

func (r resultMsg) String() string {
	if r.recvTime == "" {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("%s %s", r.content,
		durationStyle.Render(r.recvTime))
}

type model struct {
	spinner         spinner.Model
	results         []resultMsg
	textarea        textarea.Model
	quitting        bool
	upStyle         lipgloss.Style
	upTagStyle      lipgloss.Style
	ygStyle         lipgloss.Style
	ygTagStyle      lipgloss.Style
	senderStyle     lipgloss.Style
	senderTagStyle  lipgloss.Style
	welcomeTagStyle lipgloss.Style
	err             error
}

func listenMsg(p *tea.Program) {
	event.On(ON_DAN_MU, event.ListenerFunc(func(e event.Event) error {
		if s, ok := e.Data()["data"]; ok {
			p.Send(s.(*DanMuMsg))
			return nil
		} else {
			return fmt.Errorf("event:%s no data", ON_DAN_MU)
		}
	}))
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "send..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(280)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	const numLastResults = 6
	s := spinner.New()
	s.Style = spinnerStyle

	ta.KeyMap.InsertNewline.SetEnabled(false)
	m := model{
		spinner:  s,
		results:  make([]resultMsg, numLastResults),
		textarea: ta,
		quitting: false,

		upStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")),
		upTagStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		ygStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500")),
		ygTagStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#90EE9E")),
		senderStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		senderTagStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
		welcomeTagStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF	")),
		err:             nil,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textarea.Blink)
}

func now2RecvTime() string {
	now := time.Now()
	// year, mon, day := now.Date()
	h, m, _ := now.Clock()
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (m model) useStyle(dmsg *DanMuMsg, tag, sender lipgloss.Style) string {
	return tag.Render(dmsg.Tag) + sender.Render(dmsg.From) + ": " + dmsg.Content
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		if pan := recover(); pan != nil {
			log.Errorf("handler packet error: %v\n%s", pan, debug.Stack())
		}
	}()

	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		dmsg  *DanMuMsg
	)

	m.textarea, tiCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			dmsg = NewDanMuMsg(
				"[UP]",
				" ikun ",
				m.textarea.Value(),
				0,
			)
			m.results = append(
				m.results[1:],
				resultMsg{
					content:  m.useStyle(dmsg, m.upTagStyle, m.upStyle),
					recvTime: now2RecvTime(),
				},
			)
			m.textarea.Reset()
		}
	case resultMsg:
		m.results = append(m.results[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case *DanMuMsg:
		dmsg = msg
		var renderStr string
		if dmsg.Type == 1 && false {
			renderStr = m.useStyle(dmsg, m.welcomeTagStyle, m.upTagStyle)
		} else {
			if dmsg.From == " 梯度上升" {
				dmsg.Tag = "[UP大号]"
				renderStr = m.useStyle(dmsg, m.upTagStyle, m.upStyle)
			} else if dmsg.From == " 要换大房子" {
				dmsg.Tag = "[勇哥]"
				renderStr = m.useStyle(dmsg, m.ygTagStyle, m.ygStyle)
			} else {
				renderStr = m.useStyle(dmsg, m.senderTagStyle, m.senderStyle)
			}
		}
		m.results = append(
			m.results[1:],
			resultMsg{
				content:  renderStr,
				recvTime: now2RecvTime(),
			},
		)
		// m.textarea.Reset()
		// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) ViewOld() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		HEADR,
		"",
		//m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "Quit!"
	} else {
		s += m.spinner.View() + " " + HEADR
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	s += "\n" + m.textarea.View()

	return appStyle.Render(s)
}
