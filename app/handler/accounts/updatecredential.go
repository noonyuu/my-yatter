package accounts

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/auth"
)

func (h *handler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account_info := auth.AccountOf(ctx) // 認証情報を取得する

	// MultipartFormを使うには事前にParseMultipartFormを実行する
	// これがないとr.MultipartFormがnilのままとなりフォームデータにアクセスできなるっぽい
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		return
	}

	// それぞれの画像ファイルを取得する
	for _, key := range []string{"avatar", "header"} {
		files, ok := r.MultipartForm.File[key]
		// ファイルがない場合はスキップ
		if !ok {
			continue
		}
		// ファイルが複数ある場合はエラー
		if len(files) != 1 {
			http.Error(w, "Allowed up to one file.", http.StatusBadRequest)
			return
		}
		// ファイルを開く
		file, err := files[0].Open()
		// ファイルが開けない場合はエラー
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close() // 関数が終了する際にファイルを閉じる

		// ファイルのデータを読み込む
		data, err := io.ReadAll(file)
		// ファイルのデータが読み込めない場合はエラー
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ct := http.DetectContentType(data)
		if !strings.HasPrefix(ct, "image/") {
			http.Error(w, "Only image files are allowed.", http.StatusBadRequest)
			return
		}

		switch key {
		case "avatar":
			// avatarの画像パスを取得して保存
			account_info.Avatar = &files[0].Filename
		case "header":
			// headerの画像パスを取得して保存
			account_info.Header = &files[0].Filename
		}
	}
	// note, display_nameを取得して保存
	if noteValues, ok := r.MultipartForm.Value["note"]; ok {
		if len(noteValues) > 0 {
			note := noteValues[0]
			account_info.Note = &note
		}
	}
	if displayNameValues, ok := r.MultipartForm.Value["display_name"]; ok {
		if len(displayNameValues) > 0 {
			displayName := displayNameValues[0]
			account_info.DisplayName = &displayName
		}
	}

	dto, err := h.au.UpdateCredentials(ctx, account_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
