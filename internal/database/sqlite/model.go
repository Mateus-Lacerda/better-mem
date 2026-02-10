package sqlite

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	ID         string `gorm:"primaryKey"`
	ExternalID string `gorm:"uniqueIndex;not null"`
	// Has Many
	LongTermMemories  []LongTermMemory  `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	ShortTermMemories []ShortTermMemory `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type LongTermMemory struct {
	ID             string                  `gorm:"primaryKey"`
	Memory         string                  `gorm:"type:text;not null"`
	ChatID         string                  `gorm:"index:idx_ltm_query"`
	AccessCount    int                     `gorm:"default:0"`
	CreatedAt      time.Time               `gorm:"index:idx_ltm_query,priority:2;sort:desc"`
	Active         bool                    `gorm:"default:true"`
	RelatedContext []RelatedContextContent `gorm:"many2many:related_content"`
}

func (c *LongTermMemory) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type ShortTermMemory struct {
	ID             string                  `gorm:"primaryKey"`
	Memory         string                  `gorm:"type:text;not null"`
	ChatID         string                  `gorm:"index:idx_stm_query"`
	AccessCount    int                     `gorm:"default:0"`
	MergeCount     int                     `gorm:"default:0"`
	Merged         bool                    `gorm:"default:false"`
	CreatedAt      time.Time               `gorm:"index:idx_stm_query,priority:2;sort:desc"`
	Active         bool                    `gorm:"default:true"`
	RelatedContext []RelatedContextContent `gorm:"many2many:related_content"`
}

func (c *ShortTermMemory) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type RelatedContextContent struct {
	ID      string `gorm:"primaryKey"`
	Context string `gorm:"type:text;not nul"`
	User    string `gorm:"type:text;not nul"`
}

func (c *RelatedContextContent) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&Chat{},
		&LongTermMemory{},
		&ShortTermMemory{},
		&RelatedContextContent{},
	); err != nil {
		return err
	}

	return db.Exec(
		fmt.Sprintf(`
			CREATE VIRTUAL TABLE IF NOT EXISTS vec_memories USING vec0(
				embedding FLOAT[%d] distance_metric=cosine,
				id TEXT,
				memory_type INTEGER,
				+vectors_json TEXT
			);`,
			config.Database.DefaultVectorSize),
	).Error

}
