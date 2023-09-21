package lyth

import (
    "fmt"
    "sync"
    "time"
)

// Token representa um token temporário
type Token struct {
    Name     string
    Expires  time.Time
    Value    string
}

// Lythiun é a estrutura de gerenciamento da biblioteca
type Lythiun struct {
    tokens map[string]Token
    mu     sync.RWMutex
}

// New cria uma nova instância da biblioteca Lythiun
func New() *Lythiun {
    return &Lythiun{
        tokens: make(map[string]Token),
    }
}

// create cria um novo token com um nome, data de expiração e valor especificados
func (l *Lythiun) create(name string, expires string, value string) error {
    expiration, err := time.Parse("01/02/2006[15:04]", expires)
    if err != nil {
        return err
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    l.tokens[name] = Token{
        Name:    name,
        Expires: expiration,
        Value:   value,
    }
    return nil
}

// exists verifica se um token com o nome especificado existe e ainda não expirou
func (l *Lythiun) exists(name string) bool {
    l.mu.RLock()
    defer l.mu.RUnlock()

    token, ok := l.tokens[name]
    if !ok {
        return false
    }

    return token.Expires.After(time.Now())
}

// read lê o valor de um token com o nome especificado
func (l *Lythiun) read(name string) (string, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()

    token, ok := l.tokens[name]
    if !ok {
        return "", fmt.Errorf("token not found")
    }

    if token.Expires.Before(time.Now()) {
        return "", fmt.Errorf("token has expired")
    }

    return token.Value, nil
}

// get retorna a data de expiração de um token com o nome especificado
func (l *Lythiun) get(name string) (time.Time, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()

    token, ok := l.tokens[name]
    if !ok {
        return time.Time{}, fmt.Errorf("token not found")
    }

    if token.Expires.Before(time.Now()) {
        return time.Time{}, fmt.Errorf("token has expired")
    }

    return token.Expires, nil
}

// r calcula automaticamente o tempo restante para a expiração de um token
func (l *Lythiun) r(name string) (time.Duration, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()

    token, ok := l.tokens[name]
    if !ok {
        return 0, fmt.Errorf("token not found")
    }

    if token.Expires.Before(time.Now()) {
        return 0, fmt.Errorf("token has expired")
    }

    remainingTime := token.Expires.Sub(time.Now())
    return remainingTime, nil
}
