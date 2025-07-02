package schedulerstorage

type SchedulerStorrage interface {
	ListEventsToNotify()
	DeleteOlderThan()
}

func New() {

}
