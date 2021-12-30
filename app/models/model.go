package models

import "goblog/pkg/types"

// BaseModel 基类
type BaseModel struct {
	ID uint64
}

func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}