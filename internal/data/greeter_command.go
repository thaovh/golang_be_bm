package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterCommandRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterCommandRepo creates a new GreeterCommandRepo
func NewGreeterCommandRepo(data *Data, logger log.Logger) biz.GreeterCommandRepo {
	return &greeterCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Save saves a Greeter using write database
func (r *greeterCommandRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	// Sử dụng write database cho write operations
	_ = r.data.GetWriteDB() // db will be used when implementing actual logic
	r.log.WithContext(ctx).Infof("Saving greeter to write database: %v", g.Hello)
	
	// TODO: Implement actual save logic with GORM
	// Example:
	// db := r.data.GetWriteDB()
	// if err := db.WithContext(ctx).Create(g).Error; err != nil {
	//     return nil, err
	// }
	
	return g, nil
}

// Update updates a Greeter using write database
func (r *greeterCommandRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	// Sử dụng write database cho write operations
	_ = r.data.GetWriteDB() // db will be used when implementing actual logic
	r.log.WithContext(ctx).Infof("Updating greeter in write database: %v", g.Hello)
	
	// TODO: Implement actual update logic with GORM
	// Example:
	// db := r.data.GetWriteDB()
	// if err := db.WithContext(ctx).Save(g).Error; err != nil {
	//     return nil, err
	// }
	
	return g, nil
}

// Delete deletes a Greeter using write database
func (r *greeterCommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// Sử dụng write database cho write operations
	_ = r.data.GetWriteDB() // db will be used when implementing actual logic
	r.log.WithContext(ctx).Infof("Deleting greeter from write database: %s", id.String())
	
	// TODO: Implement actual delete logic with GORM
	// Example:
	// db := r.data.GetWriteDB()
	// return db.WithContext(ctx).Delete(&biz.Greeter{}, "id = ?", id).Error
	
	return nil
}

