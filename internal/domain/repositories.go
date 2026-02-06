package domain

// Definimos los métodos para la persistencia de cada tabla

type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	LoginUser(p *User) (*User, error)
}

type PostRepository interface {
	Create(post *Post) error
	Update(post *Post) error
	Delete(id int64, author string) error
	FindByID(id int64) (*Post, error)
	FindByAuthor(author string) ([]*Post, error)
}

type ProfileRepository interface {
	Create(profile *Profile) error
	Update(profile *Profile) error
	FindByUsername(username string) (*Profile, error)
}

type UserFollowRepository interface {
	Create(follow *UserFollow) error
	Delete(follow *UserFollow) error
	GetFollowers(username string) ([]string, error)
	GetFollowing(username string) ([]string, error)
}
