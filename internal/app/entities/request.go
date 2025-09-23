package entities

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,username"`
	Password  string `json:"password" binding:"required,password"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,password"`
}

type DownloadRequest struct {
	Path string `json:"path" binding:"required"`
}
