package beautiful

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/fatih/color"
)

// JSONNode representa um nó na árvore JSON
type JSONNode struct {
	Key        string
	Value      interface{}
	Type       string
	IsExpanded bool
	Children   []*JSONNode
	Parent     *JSONNode
	Level      int
}

// JSONExplorer fornece funcionalidade para explorar JSON interativamente
type JSONExplorer struct {
	root     *JSONNode
	cursor   *JSONNode
	output   *Output
	terminal *Terminal
}

// Terminal gerencia operações do terminal
type Terminal struct {
	originalState *TerminalState
}

// TerminalState armazena o estado original do terminal
type TerminalState struct {
	termios syscall.Termios
}

// NewJSONExplorer cria uma nova instância do explorador JSON
func NewJSONExplorer(output *Output) *JSONExplorer {
	return &JSONExplorer{
		output:   output,
		terminal: &Terminal{},
	}
}

// ExploreJSON inicia a exploração interativa de um JSON
func (je *JSONExplorer) ExploreJSON(data []byte) error {
	// Parse do JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	// Construir árvore de nós
	je.root = je.buildNode("root", jsonData, 0)
	je.cursor = je.root

	// Configurar terminal
	if err := je.terminal.setRawMode(); err != nil {
		return fmt.Errorf("erro ao configurar terminal: %v", err)
	}
	defer je.terminal.restoreMode()

	// Iniciar loop de navegação
	return je.navigationLoop()
}

// buildNode constrói um nó da árvore JSON
func (je *JSONExplorer) buildNode(key string, value interface{}, level int) *JSONNode {
	node := &JSONNode{
		Key:        key,
		Value:      value,
		Level:      level,
		IsExpanded: false,
		Children:   []*JSONNode{},
	}

	switch v := value.(type) {
	case map[string]interface{}:
		node.Type = "object"
		for k, val := range v {
			child := je.buildNode(k, val, level+1)
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case []interface{}:
		node.Type = "array"
		for i, val := range v {
			child := je.buildNode(strconv.Itoa(i), val, level+1)
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = "string"
		node.Value = fmt.Sprintf("\"%s\"", v)
	case float64:
		node.Type = "number"
		node.Value = v
	case bool:
		node.Type = "boolean"
		node.Value = v
	case nil:
		node.Type = "null"
		node.Value = "null"
	default:
		node.Type = "unknown"
		node.Value = fmt.Sprintf("%v", v)
	}

	return node
}

// navigationLoop gerencia o loop principal de navegação
func (je *JSONExplorer) navigationLoop() error {
	for {
		// Limpar tela
		je.clearScreen()

		// Renderizar árvore
		je.renderTree()

		// Ler input do usuário
		key := je.readKey()

		switch key {
		case "q", "Q":
			return nil // Sair
		case "up", "k":
			je.moveUp()
		case "down", "j":
			je.moveDown()
		case "left", "h":
			je.moveLeft()
		case "right", "l":
			je.moveRight()
		case "enter", " ":
			je.toggleExpand()
		case "home":
			je.goToRoot()
		case "end":
			je.goToEnd()
		}
	}
}

// renderTree renderiza a árvore JSON na tela
func (je *JSONExplorer) renderTree() {
	je.output.PrintHeader("Explorador JSON - Use as setas para navegar, Enter para expandir/colapsar, Q para sair")

	// Obter altura da tela
	_, height := je.getTerminalSize()
	availableLines := height - 5 // Reservar espaço para cabeçalho e instruções

	// Renderizar nós visíveis
	visibleNodes := je.getVisibleNodes(availableLines)

	for i, node := range visibleNodes {
		isCursor := (node == je.cursor)
		je.renderNode(node, isCursor)

		// Parar se atingiu o limite da tela
		if i >= availableLines-1 {
			break
		}
	}

	// Mostrar instruções
	je.showInstructions()
}

// getVisibleNodes retorna os nós visíveis na tela
func (je *JSONExplorer) getVisibleNodes(maxLines int) []*JSONNode {
	var visible []*JSONNode
	je.collectVisibleNodes(je.root, &visible, maxLines)
	return visible
}

// collectVisibleNodes coleta nós visíveis recursivamente
func (je *JSONExplorer) collectVisibleNodes(node *JSONNode, visible *[]*JSONNode, maxLines int) {
	if len(*visible) >= maxLines {
		return
	}

	*visible = append(*visible, node)

	if node.IsExpanded {
		for _, child := range node.Children {
			je.collectVisibleNodes(child, visible, maxLines)
		}
	}
}

// renderNode renderiza um nó individual
func (je *JSONExplorer) renderNode(node *JSONNode, isCursor bool) {
	indent := strings.Repeat("  ", node.Level)

	var prefix string
	if node.Type == "object" || node.Type == "array" {
		if node.IsExpanded {
			prefix = "▼"
		} else {
			prefix = "▶"
		}
	} else {
		prefix = "•"
	}

	// Construir linha
	line := fmt.Sprintf("%s%s %s", indent, prefix, je.formatNode(node))

	// Aplicar cores baseadas no tipo e cursor
	if isCursor {
		cursorColor := color.New(color.BgCyan, color.FgBlack, color.Bold)
		fmt.Println(cursorColor.Sprint(line))
	} else {
		je.colorizeNode(node, line)
	}
}

// formatNode formata um nó para exibição
func (je *JSONExplorer) formatNode(node *JSONNode) string {
	if node.Key == "root" {
		return "JSON Root"
	}

	switch node.Type {
	case "object":
		return fmt.Sprintf("%s: {object} (%d items)", node.Key, len(node.Children))
	case "array":
		return fmt.Sprintf("%s: [array] (%d items)", node.Key, len(node.Children))
	case "string":
		return fmt.Sprintf("%s: %s", node.Key, node.Value)
	case "number":
		return fmt.Sprintf("%s: %v", node.Key, node.Value)
	case "boolean":
		return fmt.Sprintf("%s: %v", node.Key, node.Value)
	case "null":
		return fmt.Sprintf("%s: null", node.Key)
	default:
		return fmt.Sprintf("%s: %v", node.Key, node.Value)
	}
}

// colorizeNode aplica cores ao nó baseado no tipo
func (je *JSONExplorer) colorizeNode(node *JSONNode, line string) {
	switch node.Type {
	case "object":
		objColor := color.New(color.FgBlue, color.Bold)
		fmt.Println(objColor.Sprint(line))
	case "array":
		arrColor := color.New(color.FgMagenta, color.Bold)
		fmt.Println(arrColor.Sprint(line))
	case "string":
		strColor := color.New(color.FgGreen)
		fmt.Println(strColor.Sprint(line))
	case "number":
		numColor := color.New(color.FgCyan)
		fmt.Println(numColor.Sprint(line))
	case "boolean":
		boolColor := color.New(color.FgYellow, color.Bold)
		fmt.Println(boolColor.Sprint(line))
	case "null":
		nullColor := color.New(color.FgRed, color.Bold)
		fmt.Println(nullColor.Sprint(line))
	default:
		fmt.Println(line)
	}
}

// showInstructions mostra as instruções de navegação
func (je *JSONExplorer) showInstructions() {
	fmt.Println()
	infoColor := color.New(color.FgCyan)
	infoColor.Println("Navegação: ↑↓ (mover) | →← (expandir/colapsar) | Enter (toggle) | Q (sair)")
}

// moveUp move o cursor para cima
func (je *JSONExplorer) moveUp() {
	visibleNodes := je.getVisibleNodes(1000) // Número grande para pegar todos
	for i, node := range visibleNodes {
		if node == je.cursor && i > 0 {
			je.cursor = visibleNodes[i-1]
			break
		}
	}
}

// moveDown move o cursor para baixo
func (je *JSONExplorer) moveDown() {
	visibleNodes := je.getVisibleNodes(1000)
	for i, node := range visibleNodes {
		if node == je.cursor && i < len(visibleNodes)-1 {
			je.cursor = visibleNodes[i+1]
			break
		}
	}
}

// moveLeft move o cursor para a esquerda (colapsar)
func (je *JSONExplorer) moveLeft() {
	if je.cursor.IsExpanded {
		je.cursor.IsExpanded = false
	} else if je.cursor.Parent != nil {
		je.cursor = je.cursor.Parent
	}
}

// moveRight move o cursor para a direita (expandir)
func (je *JSONExplorer) moveRight() {
	if !je.cursor.IsExpanded && (je.cursor.Type == "object" || je.cursor.Type == "array") {
		je.cursor.IsExpanded = true
	} else if je.cursor.IsExpanded && len(je.cursor.Children) > 0 {
		je.cursor = je.cursor.Children[0]
	}
}

// toggleExpand alterna a expansão do nó atual
func (je *JSONExplorer) toggleExpand() {
	if je.cursor.Type == "object" || je.cursor.Type == "array" {
		je.cursor.IsExpanded = !je.cursor.IsExpanded
	}
}

// goToRoot vai para o nó raiz
func (je *JSONExplorer) goToRoot() {
	je.cursor = je.root
}

// goToEnd vai para o último nó visível
func (je *JSONExplorer) goToEnd() {
	visibleNodes := je.getVisibleNodes(1000)
	if len(visibleNodes) > 0 {
		je.cursor = visibleNodes[len(visibleNodes)-1]
	}
}

// readKey lê uma tecla do teclado
func (je *JSONExplorer) readKey() string {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return "q"
	}

	// Verificar sequências especiais
	if char == 27 { // ESC
		next, _, _ := reader.ReadRune()
		if next == 91 { // [
			third, _, _ := reader.ReadRune()
			switch third {
			case 65:
				return "up"
			case 66:
				return "down"
			case 67:
				return "right"
			case 68:
				return "left"
			case 72:
				return "home"
			case 70:
				return "end"
			}
		}
	}

	// Verificar teclas especiais
	switch char {
	case 13:
		return "enter"
	case 32:
		return " "
	case 113, 81: // q, Q
		return "q"
	case 106, 74: // j, J
		return "down"
	case 107, 75: // k, K
		return "up"
	case 104, 72: // h, H
		return "left"
	case 108, 76: // l, L
		return "right"
	}

	return string(char)
}

// clearScreen limpa a tela
func (je *JSONExplorer) clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// getTerminalSize obtém o tamanho do terminal
func (je *JSONExplorer) getTerminalSize() (width, height int) {
	// Implementação básica - retorna valores padrão
	return 80, 24
}

// setRawMode configura o terminal em modo raw
func (t *Terminal) setRawMode() error {
	if runtime.GOOS == "windows" {
		return nil // Windows não suporta termios
	}

	fd := int(os.Stdin.Fd())
	termios, err := t.getTermios(fd)
	if err != nil {
		return err
	}

	t.originalState = &TerminalState{termios: *termios}

	// Configurar modo raw
	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	termios.Iflag &^= syscall.IXON | syscall.IXOFF | syscall.ICRNL
	termios.Cflag &^= syscall.CSIZE | syscall.PARENB
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	return t.setTermios(fd, termios)
}

// restoreMode restaura o modo original do terminal
func (t *Terminal) restoreMode() error {
	if t.originalState == nil || runtime.GOOS == "windows" {
		return nil
	}

	fd := int(os.Stdin.Fd())
	return t.setTermios(fd, &t.originalState.termios)
}

// getTermios obtém as configurações do terminal
func (t *Terminal) getTermios(fd int) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return &termios, nil
}

// setTermios define as configurações do terminal
func (t *Terminal) setTermios(fd int, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
