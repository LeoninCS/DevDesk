package CodeShare

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cs *CodeShare) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Author   string `json:"author"`
		Language string `json:"language"`
		Content  string `json:"content"`
		TTL      int64  `json:"ttl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.TTL == 0 {
		req.TTL = 3600
	}

	code := cs.Upload(req.Author, req.Language, req.Content, req.TTL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "success",
		"hash":    code.Hash,
		"url":     "/code/" + code.Hash,
	})
}

func (cs *CodeShare) GetHandler(w http.ResponseWriter, r *http.Request) {
	hash := strings.TrimPrefix(r.URL.Path, "/code/")
	code, ok := cs.Get(hash)
	if !ok {
		http.Error(w, "Code not found or expired", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(code)
}

func NewHandler(cs *CodeShare) http.Handler {
	mux := http.NewServeMux()
	// 上传代码
	mux.HandleFunc("/upload", cs.UploadHandler)
	// 获取代码
	mux.HandleFunc("/code/", cs.GetHandler)

	return mux
}
