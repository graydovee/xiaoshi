package chatgpt

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Prompt struct {
	Path        string            `json:"-" yaml:"-"`
	BasePrompt  string            `json:"basePrompt" yaml:"basePrompt"`
	Characters  map[string]string `json:"characters" yaml:"characters"`
	DefaultRole string            `json:"defaultRole" yaml:"defaultRole"`
}

func MustLoadRole(path string, data []byte) *Prompt {
	var p Prompt
	err := yaml.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}
	p.Path = path
	return &p
}

func (p *Prompt) SwitchRole(role string) {
	p.DefaultRole = role

	bytes, err := yaml.Marshal(p)
	if err == nil {
		_ = os.WriteFile(p.Path, bytes, 0644)
	}
}

func (p *Prompt) SetRolePrompt(role, prompt string) {
	if p.Characters == nil {
		p.Characters = make(map[string]string)
	}
	p.Characters[role] = prompt

	bytes, err := yaml.Marshal(p)
	if err == nil {
		_ = os.WriteFile(p.Path, bytes, 0644)
	}
}

func (p *Prompt) DeleteRolePrompt(role string) {
	delete(p.Characters, role)

	bytes, err := yaml.Marshal(p)
	if err == nil {
		_ = os.WriteFile(p.Path, bytes, 0644)
	}
}

func (p *Prompt) GetRolePrompt(role string) ([]string, bool) {
	if p == nil || p.Characters == nil {
		return nil, false
	}
	if prompt, ok := p.Characters[role]; ok {
		return []string{p.BasePrompt, prompt}, true
	}
	return nil, false
}

func (p *Prompt) GetPrompt() []string {
	return []string{p.BasePrompt, p.Characters[p.DefaultRole]}
}
