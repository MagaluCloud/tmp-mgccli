package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

//go:embed translations/*.json
var translationsFS embed.FS

// Locale representa um idioma/localização
type Locale struct {
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	NativeName   string            `json:"native_name"`
	Translations map[string]string `json:"translations"`
}

// Manager gerencia as traduções e idiomas
type Manager struct {
	locales     map[string]*Locale
	current     *Locale
	defaultLang string
	mutex       sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

func Init18n(lang string) *Manager {
	once.Do(func() {
		instance = &Manager{
			locales:     make(map[string]*Locale),
			defaultLang: lang,
		}
		if err := instance.LoadLocales(); err != nil {
			fmt.Printf("Warning: failed to load translations: %v\n", err)
		}
	})
	return instance
}

func GetInstance() *Manager {
	return instance
}

func (m *Manager) LoadLocales() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	files, err := translationsFS.ReadDir("translations")
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos de tradução: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			locale, err := m.loadLocaleFile("translations/" + file.Name())
			if err != nil {
				return fmt.Errorf("erro ao carregar %s: %w", file.Name(), err)
			}
			m.locales[locale.Code] = locale
		}
	}

	if len(m.locales) == 0 {
		return fmt.Errorf("nenhum arquivo de tradução encontrado")
	}

	m.setCurrentLocale(m.detectLanguage())

	return nil
}

func (m *Manager) loadLocaleFile(filepath string) (*Locale, error) {
	data, err := translationsFS.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var locale Locale
	if err := json.Unmarshal(data, &locale); err != nil {
		return nil, err
	}

	return &locale, nil
}

func (m *Manager) detectLanguage() string {
	if m.defaultLang != "" {
		return m.defaultLang
	}

	if lang := os.Getenv("CLI_LANG"); lang != "" {
		if m.isValidLocale(lang) {
			return lang
		}
	}

	if lang := os.Getenv("LANG"); lang != "" {
		langCode := strings.Split(lang, ".")[0]
		langCode = strings.Replace(langCode, "_", "-", 1)

		if m.isValidLocale(langCode) {
			return langCode
		}

		mainLang := strings.Split(langCode, "-")[0]
		if m.isValidLocale(mainLang) {
			return mainLang
		}
	}

	if lang := os.Getenv("LC_ALL"); lang != "" {
		langCode := strings.Split(lang, ".")[0]
		langCode = strings.Replace(langCode, "_", "-", 1)

		if m.isValidLocale(langCode) {
			return langCode
		}
	}

	return "pt-BR"
}

func (m *Manager) isValidLocale(code string) bool {
	_, exists := m.locales[code]
	return exists
}

func (m *Manager) setCurrentLocale(code string) {
	if locale, exists := m.locales[code]; exists {
		m.current = locale
	} else {
		for _, locale := range m.locales {
			m.current = locale
			break
		}
	}
}

func (m *Manager) SetLanguage(code string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.isValidLocale(code) {
		return fmt.Errorf("idioma não suportado: %s", code)
	}

	m.setCurrentLocale(code)
	return nil
}

func (m *Manager) GetLanguage() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.current == nil {
		return m.defaultLang
	}
	return m.current.Code
}

func (m *Manager) T(key string, args ...interface{}) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.current == nil {
		return key
	}

	translation, exists := m.current.Translations[key]
	if !exists {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(translation, args...)
	}

	return translation
}

func (m *Manager) GetAvailableLanguages() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var languages []string
	for code := range m.locales {
		languages = append(languages, code)
	}
	return languages
}

func (m *Manager) GetLanguageInfo(code string) (*Locale, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	locale, exists := m.locales[code]
	if !exists {
		return nil, fmt.Errorf("idioma não encontrado: %s", code)
	}

	return locale, nil
}

func (m *Manager) SetupCobraI18n(cmd *cobra.Command) {
	cmd.PersistentFlags().String("lang", "", "Idioma da interface (ex: pt-BR, en-US)")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if lang, _ := cmd.Flags().GetString("lang"); lang != "" {
			return m.SetLanguage(lang)
		}
		return nil
	}
}
