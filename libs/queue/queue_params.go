package queue

type QueueParams struct {
	name        string
	qosCount    int
	autoAck     bool
	kind        string
	durable     bool
	internal    bool
	autoDelete  bool
	exclusive   bool
	noWait      bool
	threadLimit int
	args        map[string]interface{}
}

func NewQueueParams(name string) QueueParams {
	q := QueueParams{
		name: name,
	}
	q.LoadDefaults()
	return q
}

func (qp *QueueParams) LoadDefaults() {
	qp.SetDurable(true)
	qp.SetAutoAck(false)
	qp.SetKind("direct")
	qp.SetInternal(false)
	qp.SetAutoDelete(false)
	qp.SetExclusive(false)
	qp.SetNoWait(false)
	qp.SetArgs(nil)
	qp.SetThreadLimit(5)
}

func (q *QueueParams) AutoAck() bool {
	return q.autoAck
}

func (q *QueueParams) SetAutoAck(autoAck bool) {
	q.autoAck = autoAck
}

func (q *QueueParams) Args() map[string]interface{} {
	return q.args
}

func (q *QueueParams) SetArgs(args map[string]interface{}) {
	q.args = args
}

func (q *QueueParams) Name() string {
	return q.name
}

func (q *QueueParams) Internal() bool {
	return q.internal
}

func (q *QueueParams) SetInternal(internal bool) {
	q.internal = internal
}

func (q *QueueParams) NoWait() bool {
	return q.noWait
}

func (q *QueueParams) SetNoWait(noWait bool) {
	q.noWait = noWait
}

func (q *QueueParams) Exclusive() bool {
	return q.exclusive
}

func (q *QueueParams) SetExclusive(exclusive bool) {
	q.exclusive = exclusive
}

func (q *QueueParams) AutoDelete() bool {
	return q.autoDelete
}

func (q *QueueParams) SetAutoDelete(autoDelete bool) {
	q.autoDelete = autoDelete
}

func (q *QueueParams) Durable() bool {
	return q.durable
}

func (q *QueueParams) SetDurable(durable bool) {
	q.durable = durable
}

func (q *QueueParams) Kind() string {
	return q.kind
}

func (q *QueueParams) SetKind(kind string) {
	q.kind = kind
}

func (q *QueueParams) ThreadLimit() int {
	return q.threadLimit
}

func (q *QueueParams) SetThreadLimit(threadLimit int) {
	q.threadLimit = threadLimit
}
