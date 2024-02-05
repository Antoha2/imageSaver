package repository

type Repository interface{}

type repositoryImpl struct{}

func NewRepository() *repositoryImpl {

	return &repositoryImpl{}
}
