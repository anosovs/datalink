package handler

import (
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/anosovs/datalink/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage sqlite.Storage
}

func Init (storage sqlite.Storage) *Handler{
	return &Handler{
		Storage: storage,
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

	Link := URL{
		Link: "",
	}
	if r.TLS != nil {
		Link.Link = "https://"
	} else {
		Link.Link = "http://"
	}
	
	Link.Link += r.Host +"/"+ uuid
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