package model

import (
	"fmt"
)

// Reference379 — представление табл. товаров из MSSQL
type Reference379 struct {
	ID           []byte  `db:"_IDRRef"`       // уникальный идентификатор
	Version      []byte  `db:"_Version"`      // версия
	Marked       BitBool `db:"_Marked"`       // пометка удаления
	PredefinedID string  `db:"_PredefinedID"` // резерв
	ParentID     []byte  `db:"_ParentIDRRef"` // ссылка на родителя
	Folder       BitBool `db:"_Folder"`       // признак папки/группы
	Code         string  `db:"_Code"`         // код номенклатуры
	Description  string  `db:"_Description"`  // наименование
}


// BitBool - bool версия MSSQL
type BitBool bool

func (bb *BitBool) Scan(src any) error {
	if src == nil {
		*bb = false
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for BitBool, got %T", src)
	}
	*bb = BitBool(len(b) > 0 && b[0] == 1)
	return nil
}
