package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *WebServer) Route() {
	s.mux.Use(middleware.Recoverer, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// The home route notifies that the API is up and running
	s.mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SOCIALAT API is up and running"))
	})
	s.mux.Get("/socket.io/", s.handleSocket())
	fs := http.FileServer(http.Dir("upload"))
	s.mux.Handle("/upload/*", http.StripPrefix("/upload/", fs))
	s.mux.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			var authRouter = apiAuth{WebServer: s}
			r.Get("/auth-method", authRouter.getAuthMethod)
			r.Post("/assertion-options", authRouter.AssertionOptions)
			r.Post("/assertion-result", authRouter.AssertionResult)
			r.Get("/check-auth-username", authRouter.CheckAuthUsername)
			r.Get("/gen-random-username", authRouter.GenRandomUsername)
			r.Post("/cancel-register", authRouter.CancelPasskeyRegister)
			r.Post("/register-start", authRouter.StartPasskeyRegister)
			r.Post("/register-finish", authRouter.FinishPasskeyRegister)
			r.Post("/register-transfer-finish", authRouter.FinishPasskeyTransferRegister)
			r.Post("/update-passkey-start", authRouter.UpdatePasskeyStart)
			r.Post("/update-passkey-finish", authRouter.UpdatePasskeyFinish)
			r.Post("/register", authRouter.register)
			r.Post("/login", authRouter.login)
		})
		r.Route("/user", func(r chi.Router) {
			var userRouter = apiUser{WebServer: s}
			r.Use(s.loggedInMiddleware)
			r.Post("/update-display-name", userRouter.updateDisplayName)
			r.Post("/update-bio", userRouter.updateBio)
			r.Post("/update-avatar", userRouter.UpdateAvatar)
			r.Post("/follow-user", userRouter.FollowUpdateUser)
		})
		r.Route("/public", func(r chi.Router) {
			var userRouter = apiUser{WebServer: s}
			r.Use(s.getLoginInfoMiddleware)
			r.Get("/get-user-by-name/{username}", userRouter.getUserByName)
			r.Get("/get-timelines", userRouter.getTimelines)
			r.Get("/get-all-posts", userRouter.getAllPosts)
			r.Get("/get-user-posts/{username}", userRouter.getUserPosts)
		})
		r.Route("/file", func(r chi.Router) {
			r.Use(s.loggedInMiddleware)
			var fileRouter = apiFileUpload{WebServer: s}
			r.Post("/upload", fileRouter.uploadFiles)
			r.Post("/upload-one", fileRouter.uploadOneFile)
			r.Get("/base64", fileRouter.getProductImagesBase64)
			r.Get("/base64-one", fileRouter.getOneImageBase64)
			r.Get("/img-base64", fileRouter.getImageBase64)
		})
		r.Route("/post", func(r chi.Router) {
			r.Use(s.loggedInMiddleware)
			var postRouter = apiPost{WebServer: s}
			r.Post("/upload-images", postRouter.uploadImages)
			r.Post("/post-with-files", postRouter.PostWithFiles)
			r.Post("/post-without-files", postRouter.PostWithoutFiles)
		})
	})
}
