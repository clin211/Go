package model

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"` // 自增长id
	CreatedBy  string `json:"created_by"`            // 创建人
	ModifiedBy string `json:"modified_by"`           // 修改人
	CreatedOn  uint32 `json:"created_on"`            // 创建时间
	ModifiedOn uint32 `json:"modified_on"`           // 修改时间
	DeletedOn  uint32 `json:"deleted_on"`            // 删除时间
	IsDel      uint8  `json:"is_del"`
}
