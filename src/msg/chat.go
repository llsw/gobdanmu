package msg

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gookit/event"
)

const (
	HEADR = `欢迎弹弹幕指导☝️ 🤓`
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
	From    string
	Content string
}

// 简单的情况就不要用opts模式了
func NewDanMuMsg(from, content string) *DanMuMsg {
	return &DanMuMsg{
		From:    from,
		Content: content,
	}
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
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

	ta.SetWidth(30)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	// vp.SetContent(HEADR)

	ta.KeyMap.InsertNewline.SetEnabled(false)
	m := model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("[up] ikun: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case *DanMuMsg:
		dmsg := msg
		from := fmt.Sprintf("%s: ", dmsg.From)
		m.messages = append(m.messages, m.senderStyle.Render(from)+dmsg.Content)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		HEADR,
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
