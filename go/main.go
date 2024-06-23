package main

import (
    "fmt"
    "math"
    "os"
    "strings"

    ctp "github.com/catppuccin/go"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/huh"
    "github.com/charmbracelet/lipgloss"
    "golang.org/x/text/cases"
    lang "golang.org/x/text/language"
)

var (
    maxWidth = 80
    red      = lipgloss.AdaptiveColor{Light: ctp.Latte.Red().Hex, Dark: ctp.Mocha.Red().Hex}
    indigo   = lipgloss.AdaptiveColor{Light: ctp.Latte.Mauve().Hex, Dark: ctp.Mocha.Mauve().Hex}
    green    = lipgloss.AdaptiveColor{Light: ctp.Latte.Green().Hex, Dark: ctp.Mocha.Green().Hex}
)

type Styles struct {
    Base,
    HeaderText,
    Status,
    StatusHeader,
    Highlight,
    ErrorHeaderText,
    Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
    s := Styles{}
    s.Base = lg.NewStyle().Padding(1, 4, 0, 1)
    s.HeaderText = lg.NewStyle().Foreground(indigo).Bold(true).Padding(0, 1, 0, 2)
    s.Status = lg.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(indigo).PaddingLeft(1).MarginTop(1)
    s.StatusHeader = lg.NewStyle().Foreground(indigo).Bold(true)
    s.Highlight = lg.NewStyle().Foreground(lipgloss.Color("212"))
    s.ErrorHeaderText = s.HeaderText.Foreground(red)
    s.Help = lg.NewStyle().Foreground(lipgloss.Color("240"))

    return &s
}

const (
    statusNormal int = iota
    stateDone
)

type Model struct {
    state  int
    lg     *lipgloss.Renderer
    styles *Styles
    form   *huh.Form
    width  int
}

func NewModel() Model {
    m := Model{width: maxWidth}
    m.lg = lipgloss.DefaultRenderer()
    m.styles = NewStyles(m.lg)

    m.form = huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().Key("isTypeScript").Title("Do you want to use TypeScript?"),
            huh.NewSelect[string]().Key("framework").Title("Choose a framework").Options(
                huh.NewOptions(
                    "Express",
                    "Hono",
                    "Nest",
                    "Next",
                    "Nuxt",
                    "SvelteKit",
                    "SolidStart",
                    "Vite + React",
                    "Vite + Vue",
                    "Vite + Svelte",
                    "Vite + Solid",
                )...,
            ),
            huh.NewSelect[string]().Key("uiFramework").Title("Choose a UI framework").Options(
                huh.NewOptions(
                    "None",
                    "shadcn/ui",
                    "Tailwind UI",
                    "Daisy UI",
                    "Material UI",
                    "Chakra",
                    "Semantic UI",
                )...,
            ),
            huh.NewSelect[string]().Key("cssFramework").Title("Choose a CSS framework").Options(
                huh.NewOptions(
                    "None",
                    "TailwindCSS",
                    "Emotion",
                    "StyledCSS",
                    "Bootstrap",
                )...,
            ),
            huh.NewSelect[string]().Key("database").Title("Choose a relational database management system").Options(
                huh.NewOptions("None", "PostgreSQL", "MySQL", "SQLite")...,
            ),
            huh.NewSelect[string]().Key("orm").Title("Choose an object-relational mapper").Options(
                huh.NewOptions("None", "Prisma", "Drizzle")...,
            ),
            huh.NewSelect[string]().Key("cloudPlatform").Title("Where will you be deploying to?").Options(
                huh.NewOptions(
                    "None",
                    "GCP Cloud Run",
                    "GCP Cloud Functions",
                    "GCP App Engine",
                    "AWS Lambda",
                    "AWS Elastic Container Service",
                    "Azure Container Service",
                    "Vercel",
                )...,
            ),
            huh.NewConfirm().Key("isDocker").Title("Do you want to use Docker?"),
        ),
    ).WithWidth(45).WithShowHelp(false).WithShowErrors(false)

    return m
}

func (m Model) Init() tea.Cmd {
    return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = int(math.Min(float64(msg.Width), float64(maxWidth)))
    case tea.KeyMsg:
        switch msg.String() {
        case "esc", "ctrl+c", "q":
            return m, tea.Quit
        }
    }

    var cmds []tea.Cmd

    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
        cmds = append(cmds, cmd)
    }

    if m.form.State == huh.StateCompleted {
        cmds = append(cmds, tea.Quit)
    }

    return m, tea.Batch(cmds...)
}

func (m Model) getFramework() (string, string) {
    framework := m.form.GetString("framework")
    return framework, cases.Title(lang.English).String(framework)
}

func (m Model) errorView() string {
    var s string
    for _, err := range m.form.Errors() {
        s += err.Error()
    }
    return s
}

func (m Model) appBoundaryView(text string) string {
    return lipgloss.PlaceHorizontal(
        m.width,
        lipgloss.Left,
        m.styles.HeaderText.Render(text),
        lipgloss.WithWhitespaceChars("/"),
        lipgloss.WithWhitespaceForeground(indigo),
    )
}

func (m Model) appErrorBoundaryView(text string) string {
    return lipgloss.PlaceHorizontal(
        m.width,
        lipgloss.Left,
        m.styles.ErrorHeaderText.Render(text),
        lipgloss.WithWhitespaceChars("/"),
        lipgloss.WithWhitespaceForeground(red),
    )
}

func (m Model) View() string {
    s := m.styles

    switch m.form.State {
    case huh.StateCompleted:
        framework, description := m.getFramework()
        framework = s.Highlight.Render(framework)
        var b strings.Builder
        fmt.Fprintf(&b, "You chose %s\n", framework)
        fmt.Fprintf(&b, "This framework is %s\n", description)
        return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
    default:
        v := strings.TrimSuffix(m.form.View(), "\n\n")
        form := m.lg.NewStyle().Margin(1, 0).Render(v)

        var status string
        {
            var (
                isTypeScript  bool
                framework     string
                uiFramework   string
                cssFramework  string
                database      string
                orm           string
                cloudPlatform string
                isDocker      string
            )

            isTypeScript = m.form.GetBool("isTypeScript")
            var language string
            if isTypeScript {
                language = "TypeScript"
            } else {
                language = "JavaScript"
            }
            language = fmt.Sprintf(
                "%s\n%s\n\n",
                s.StatusHeader.Render("Language"),
                language,
            )

            if m.form.GetString("framework") != "" {
                framework = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("Framework"),
                    m.form.GetString("framework"),
                )
            }

            if m.form.GetString("uiFramework") != "" {
                uiFramework = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("UI Framework"),
                    m.form.GetString("uiFramework"),
                )
            }

            if m.form.GetString("cssFramework") != "" {
                cssFramework = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("CSS Framework"),
                    m.form.GetString("cssFramework"),
                )
            }

            if m.form.GetString("database") != "" {
                database = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("Database"),
                    m.form.GetString("database"),
                )
            }

            if m.form.GetString("orm") != "" {
                orm = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("ORM"),
                    m.form.GetString("orm"),
                )
            }

            if m.form.GetString("cloudPlatform") != "" {
                cloudPlatform = fmt.Sprintf(
                    "%s\n%s\n\n",
                    s.StatusHeader.Render("Cloud Platform"),
                    m.form.GetString("cloudPlatform"),
                )
            }

            if m.form.GetBool("isDocker") {
                isDocker = "Yes"
            } else {
                isDocker = "No"
            }
            isDocker = fmt.Sprintf(
                "%s\n%s\n\n",
                s.StatusHeader.Render("Use Docker"),
                isDocker,
            )

            const statusWidth = 28
            statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
            status = s.Status.Height(lipgloss.Height(form)).Width(statusWidth).MarginLeft(statusMarginLeft).Render(
                s.StatusHeader.Render("Tech Stack") + "\n\n" +
                    language +
                    framework +
                    uiFramework +
                    cssFramework +
                    database +
                    orm +
                    cloudPlatform +
                    isDocker,
            )
        }

        errors := m.form.Errors()
        header := m.appBoundaryView("Stack Builder")
        if len(errors) > 0 {
            header = m.appErrorBoundaryView(m.errorView())
        }
        body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

        footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
        if len(errors) > 0 {
            footer = m.appErrorBoundaryView("")
        }

        return s.Base.Render(header + "\n" + body + "\n\n" + footer)
    }
}

func main() {
    if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}
