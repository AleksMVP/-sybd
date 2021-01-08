package thread

import (
	"github.com/AleksMVP/sybd/models"
)

type IThreadUseCase interface {
	CreateThread(thread models.Thread) (t models.Thread, e error)
	GetThread(slugOrId string) (t models.Thread, e error)
	GetThreads(slug, limit, since, desc string) (ts models.Threads, e error)
	EditThread(slugOrId string, thread models.Thread) (t models.Thread, e error)
}