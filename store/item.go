package store

type Item struct {
	Id          int     `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	OrderId     string  `json:"order_id" db:"order_id"`
	ChrtId      int64   `json:"chrt_id" db:"chrt_id"`
	TrackNumber string  `json:"track_number" db:"track_number"`
	Price       float64 `json:"price" db,:"price"`
	Rid         string  `json:"rid" db:"rid"`
	Name        string  `json:"name" db:"name"`
	Sale        int     `json:"sale" db:"sale"`
	Size        string  `json:"size" db:"size"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
	NmId        int     `json:"nm_id" db:"nm_id"`
	Brand       string  `json:"brand" db:"brand"`
	Status      int     `json:"status" db:"status"`
}
