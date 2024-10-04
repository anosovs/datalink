package handler

import (
	"net/http"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/anosovs/datalink/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage storage.Storage
	DeleteAfter int
	EnableHttps bool
}

func Init (storage storage.Storage, deleteAfter int, enableHttps bool) *Handler{
	return &Handler{
		Storage: storage,
		DeleteAfter: deleteAfter,
		EnableHttps: enableHttps,
	}
}


func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
    tmpl, err := template.ParseFiles(fp)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	if err := tmpl.Execute(w, ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

type URL struct {
	Link string
	DeleteAfter string
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	message := r.Form.Get("message")
	count, err := strconv.Atoi(r.Form.Get("count"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, err := h.Storage.SaveMsg(message, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	availableUntil := time.Now().Add(time.Hour * 24 * time.Duration(h.DeleteAfter))

	Link := URL{
		Link: "",
		DeleteAfter: availableUntil.Format("2006-01-02 15:04:05"),
	}
	
	if h.EnableHttps {
		Link.Link = "https://"
	} else {
		Link.Link = "http://"
	}
	
	Link.Link += r.Host +"/show/"+ uuid
	fp := path.Join("templates", "created.html")
    tmpl, err := template.ParseFiles(fp)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	if err := tmpl.Execute(w, Link); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


type Message struct {
	Message string
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	err := h.Storage.DeleteOldMessages(h.DeleteAfter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	msg, cnt, err := h.Storage.GetMsg(uuid)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
        return
	}
	if cnt <= 1 {
		h.Storage.DeleteMsg(uuid)
	} else {
		h.Storage.DecreaseCountMsg(uuid)
	}
	data := Message{
		Message: msg,
	}

	fp := path.Join("templates", "show.html")
    tmpl, err := template.ParseFiles(fp)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}