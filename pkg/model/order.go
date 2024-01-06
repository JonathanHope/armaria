package armaria

// Order is the field results are ordered on.
type Order string

const (
	OrderModified Order = "modified"
	OrderName     Order = "name"
	OrderManual   Order = "manual"
)
