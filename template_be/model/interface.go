package model

import "context"

/*
	interface: RDbService
	description: interface of service layer for test
*/
type RDbService interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}

/*
	interface: RDbDAL
	description: interface of data access layer for test
*/
type RDbDAL interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}
