package beautiful

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Output struct {
	rawMode bool
	data    interface{}
}

func NewOutput(rawMode bool) *Output {

	return &Output{
		rawMode: rawMode,
	}
}

func (bo *Output) PrintData(data interface{}) {
	bo.data = data

	if bo.rawMode {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonData))
		return
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}

	expFlag := os.Getenv("EXPLORE_JSON") == "1"
	if expFlag {
		explorer := NewJSONExplorer(bo)
		if err := explorer.ExploreJSON(jsonData); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	coloredJSON := bo.colorizeJSON(string(jsonData))
	fmt.Println(coloredJSON)
}

func (bo *Output) PrintJSON(data interface{}) error {
	if bo.rawMode {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	coloredJSON := bo.colorizeJSON(string(jsonData))
	fmt.Println(coloredJSON)
	return nil
}

func (bo *Output) PrintSuccess(message string) {
	if bo.rawMode {
		fmt.Println(message)
		return
	}

	successColor := color.New(color.FgGreen, color.Bold)
	successColor.Printf("%s\n", message)
}

func (bo *Output) PrintError(message string) {
	if bo.rawMode {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
		return
	}

	errorColor := color.New(color.FgRed, color.Bold)
	errorColor.Printf("Error: %s\n", message)
}

func (bo *Output) PrintWarning(message string) {
	if bo.rawMode {
		fmt.Printf("Warning: %s\n", message)
		return
	}

	warningColor := color.New(color.FgYellow, color.Bold)
	warningColor.Printf("%s\n", message)
}

func (bo *Output) PrintInfo(message string) {
	if bo.rawMode {
		fmt.Println(message)
		return
	}

	infoColor := color.New(color.FgCyan, color.Bold)
	infoColor.Printf("%s\n", message)
}

func (bo *Output) PrintTable(headers []string, rows [][]string) {
	if bo.rawMode {
		fmt.Println(strings.Join(headers, "\t"))
		for _, row := range rows {
			fmt.Println(strings.Join(row, "\t"))
		}
		return
	}

	headerColor := color.New(color.FgMagenta, color.Bold)
	rowColor := color.New(color.FgWhite)

	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	headerColor.Print("┌")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┬")
		}
	}
	headerColor.Println("┐")

	headerColor.Print("│")
	for i, header := range headers {
		headerColor.Printf(" %-*s │", widths[i], header)
	}
	headerColor.Println()

	headerColor.Print("├")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┼")
		}
	}
	headerColor.Println("┤")

	for _, row := range rows {
		rowColor.Print("│")
		for i, cell := range row {
			if i < len(widths) {
				rowColor.Printf(" %-*s │", widths[i], cell)
			}
		}
		rowColor.Println()
	}

	headerColor.Print("└")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┴")
		}
	}
	headerColor.Println("┘")
}

func (bo *Output) PrintList(title string, items []string) {
	if bo.rawMode {
		fmt.Println(title)
		for _, item := range items {
			fmt.Printf("- %s\n", item)
		}
		return
	}

	titleColor := color.New(color.FgBlue, color.Bold)
	titleColor.Printf("%s:\n", title)

	itemColor := color.New(color.FgCyan)
	for i, item := range items {
		itemColor.Printf("  %d. %s\n", i+1, item)
	}
}

func (bo *Output) colorizeJSON(jsonStr string) string {
	if bo.rawMode {
		return jsonStr
	}

	braceColor := color.New(color.FgWhite, color.Bold)
	keyColor := color.New(color.FgYellow, color.Bold)
	stringColor := color.New(color.FgGreen)
	numberColor := color.New(color.FgCyan)
	booleanColor := color.New(color.FgMagenta, color.Bold)
	nullColor := color.New(color.FgRed, color.Bold)

	lines := strings.Split(jsonStr, "\n")
	var result []string

	for _, line := range lines {
		coloredLine := bo.colorizeJSONLine(line, braceColor, keyColor, stringColor, numberColor, booleanColor, nullColor)
		result = append(result, coloredLine)
	}

	return strings.Join(result, "\n")
}

func (bo *Output) colorizeJSONLine(line string, braceColor, keyColor, stringColor, numberColor, booleanColor, nullColor *color.Color) string {
	if bo.rawMode {
		return line
	}

	var result strings.Builder
	inString := false
	escapeNext := false
	afterColon := false
	i := 0

	for i < len(line) {
		char := rune(line[i])

		if escapeNext {
			result.WriteRune(char)
			escapeNext = false
			i++
			continue
		}

		if char == '\\' {
			escapeNext = true
			result.WriteRune(char)
			i++
			continue
		}

		if char == '"' {
			if !inString {
				inString = true
				if afterColon {
					result.WriteString(stringColor.Sprint(string(char)))
				} else {
					result.WriteString(keyColor.Sprint(string(char)))
				}
			} else {
				inString = false
				if afterColon {
					result.WriteString(stringColor.Sprint(string(char)))
				} else {
					result.WriteString(keyColor.Sprint(string(char)))
				}
			}
			i++
			continue
		}

		if inString {
			if afterColon {
				result.WriteString(stringColor.Sprint(string(char)))
			} else {
				result.WriteString(keyColor.Sprint(string(char)))
			}
			i++
			continue
		}

		switch char {
		case '{', '}', '[', ']':
			result.WriteString(braceColor.Sprint(string(char)))
			i++
		case ',':
			result.WriteString(braceColor.Sprint(string(char)))
			afterColon = false
			i++
		case ':':
			result.WriteString(braceColor.Sprint(string(char)))
			afterColon = true
			i++
		case ' ', '\t':
			result.WriteRune(char)
			i++
		default:
			if afterColon {
				if char == 't' && i+3 < len(line) && line[i:i+4] == "true" {
					result.WriteString(booleanColor.Sprint("true"))
					i += 4
				} else if char == 'f' && i+4 < len(line) && line[i:i+5] == "false" {
					result.WriteString(booleanColor.Sprint("false"))
					i += 5
				} else if char == 'n' && i+3 < len(line) && line[i:i+4] == "null" {
					result.WriteString(nullColor.Sprint("null"))
					i += 4
				} else if (char >= '0' && char <= '9') || char == '-' {
					result.WriteString(numberColor.Sprint(string(char)))
					i++
				} else {
					result.WriteRune(char)
					i++
				}
			} else {
				result.WriteRune(char)
				i++
			}
		}
	}

	return result.String()
}

func (bo *Output) PrintProgress(current, total int, message string) {
	if bo.rawMode {
		fmt.Printf("%s: %d/%d\n", message, current, total)
		return
	}

	progressColor := color.New(color.FgBlue, color.Bold)
	progressColor.Printf("%s: %d/%d\n", message, current, total)
}

func (bo *Output) PrintHeader(title string) {
	if bo.rawMode {
		fmt.Printf("\n=== %s ===\n", title)
		return
	}

	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Printf("\n%s\n", title)
	headerColor.Println(strings.Repeat("─", len(title)+4))
}
