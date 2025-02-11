package middleware

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type LogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Method      string    `json:"method"`
    Path        string    `json:"path"`
    StatusCode  int       `json:"status_code"`
    Duration    string    `json:"duration"`
    IPAddress   string    `json:"ip_address"`
    UserAgent   string    `json:"user_agent"`
    RequestID   string    `json:"request_id"`
    SessionID   string    `json:"session_id,omitempty"`
    Error       string    `json:"error,omitempty"`
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

type Logger struct {
    file *os.File
}

func NewLogger(logPath string) (*Logger, error) {
    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }
    return &Logger{file: file}, nil
}

func (l *Logger) LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:    http.StatusOK,
        }

        sessionID := ""
        if cookie, err := r.Cookie("session_id"); err == nil {
            sessionID = cookie.Value
        }

        next.ServeHTTP(rw, r)

        entry := LogEntry{
            Timestamp:   time.Now(),
            Method:      r.Method,
            Path:       r.URL.Path,
            StatusCode: rw.statusCode,
            Duration:   time.Since(start).String(),
            IPAddress:  r.RemoteAddr,
            UserAgent:  r.UserAgent(),
            RequestID:  r.Header.Get("X-Request-ID"),
            SessionID: sessionID,
        }

        jsonEntry, err := json.Marshal(entry)
        if err != nil {
            log.Printf("Error marshaling log entry: %v", err)
            return
        }

        if _, err := l.file.Write(append(jsonEntry, '\n')); err != nil {
            log.Printf("Error writing to log file: %v", err)
        }

        log.Printf("%s %s %d %v %s",
            entry.Method,
            entry.Path,
            entry.StatusCode,
            entry.Duration,
            entry.IPAddress,
        )
    })
}

func (l *Logger) Close() error {
    return l.file.Close()
}

func (l *Logger) LogEvent(eventType string, data interface{}) error {
    event := struct {
        Timestamp time.Time   `json:"timestamp"`
        Type      string      `json:"type"`
        Data      interface{} `json:"data"`
    }{
        Timestamp: time.Now(),
        Type:      eventType,
        Data:      data,
    }

    jsonEvent, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }

    if _, err := l.file.Write(append(jsonEvent, '\n')); err != nil {
        return fmt.Errorf("failed to write event: %w", err)
    }

    return nil
}