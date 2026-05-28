package kelurahan

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type MasterKelurahan struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	IDKecamatan   int32  `json:"id_kecamatan"`
	NameKecamatan string `json:"name_kecamatan"`
	IDCity        int32  `json:"id_city"`
	NameCity      string `json:"name_city"`
	IDProvince    int32  `json:"id_province"`
	NameProvince  string `json:"name_province"`
}

type IKelurahanRepository interface {
	GetKelurahan(ctx context.Context, skip, take int32, filter, sort string) ([]MasterKelurahan, int32, error)
	GetKelurahanById(ctx context.Context, id uuid.UUID) (MasterKelurahan, error)
	GetKelurahanByKecamatanId(ctx context.Context, kecamatanId int32) ([]MasterKelurahan, error)
	GetAllKelurahans(ctx context.Context) ([]MasterKelurahan, error)
}

type repository struct{ db *sql.DB }

func NewKelurahanRepository(db *sql.DB) IKelurahanRepository { return &repository{db: db} }

func (r *repository) GetKelurahan(ctx context.Context, skip, take int32, filter, sort string) ([]MasterKelurahan, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetKelurahanById(ctx context.Context, id uuid.UUID) (MasterKelurahan, error) {
	return MasterKelurahan{}, nil
}
func (r *repository) GetKelurahanByKecamatanId(ctx context.Context, kecamatanId int32) ([]MasterKelurahan, error) {
	return nil, nil
}
func (r *repository) GetAllKelurahans(ctx context.Context) ([]MasterKelurahan, error) {
	return nil, nil
}
