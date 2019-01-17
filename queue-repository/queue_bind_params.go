package queue_repository

type QueueBindParams struct {
	name     string
	key      string
	exchange string
	noWait   bool
	args     map[string]interface{}
}

func (q *QueueBindParams) Name() string {
	return q.name
}

func (q *QueueBindParams) Key() string {
	return q.key
}

func (q *QueueBindParams) Exchange() string {
	return q.exchange
}

func (q *QueueBindParams) Args() map[string]interface{} {
	return q.args
}

func (q *QueueBindParams) SetArgs(args map[string]interface{}) {
	q.args = args
}

func (q *QueueBindParams) NoWait() bool {
	return q.noWait
}

func (q *QueueBindParams) SetNoWait(noWait bool) {
	q.noWait = noWait
}

func NewQueueBindParams(name, key, exchange string) QueueBindParams {
	return QueueBindParams{
		name:     name,
		key:      key,
		exchange: exchange,
		noWait:   false,
		args:     nil,
	}
}
