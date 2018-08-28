package main

type Controller struct {
	ID          int `json:"id"`
	Num_Leds    int `json:"num_leds"`
	Start_Index int `json:"start_index"`
}

type DomeConfig struct {
	ControllerList []Controller `json:"Controllers"`
	LEDs           [][]float32  `json:"led_list"`
}

type LEDWithId struct {
	ID  int       `json:"id"`
	LED []float32 `json:"led"`
}
