package gotask



type _Stat struct {
	QL     int     `json:"q_l,omitempty"`
	CurDur float64 `json:"cur_dur,omitempty"`
	MaxQ   int     `json:"max_q,omitempty"`
	MaxDur float64 `json:"max_dur,omitempty"`
}
