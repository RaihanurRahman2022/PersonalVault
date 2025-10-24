package entities

// LoginRequest represents user login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required,username" example:"john_doe"`
	Password string `json:"password" binding:"required,password" example:"securePassword123"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,username" example:"john_doe"`
	Password  string `json:"password" binding:"required,password" example:"securePassword123"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50" example:"John"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50" example:"Doe"`
}

// DownloadRequest represents file download request
type DownloadRequest struct {
	Path string `json:"path" binding:"required" example:"/documents/file.pdf"`
}

// CreateFolderRequest represents folder creation request
type CreateFolderRequest struct {
	Path string `json:"path" binding:"required" example:"/documents/new_folder"`
}
