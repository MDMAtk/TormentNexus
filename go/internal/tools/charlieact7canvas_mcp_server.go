        {
          "width": 800,
          "height": 600,
          "bg": "#FFFFFF",
          "elements": [
            {"type": "rect", "x": 10, "y": 10, "w": 50, "h": 50, "color": "#FF0000"}
          ]
        }
        type Canvas struct {
            Width    int       `json:"width"`
            Height   int       `json:"height"`
            Bg       string    `json:"bg"`
            Elements []Element `json:"elements"`
        }
        type Element struct {
            Type  string `json:"type"` // "rect" or "circle"
            X     int    `json:"x"`
            Y     int    `json:"y"`
            W     int    `json:"w,omitempty"`
            H     int    `json:"h,omitempty"`
            R     int    `json:"r,omitempty"`
            Color string `json:"color"`
        }
    func drawCircle(img *image.RGBA, x0, y0, r int, c color.Color) {
        f := 1 - r
        dx := 1
        dy := -2 * r
        x := 0
        y := r
        for x < y {
            img.Set(x0+x, y0+y, c)
            img.Set(x0-x, y0+y, c)
            img.Set(x0+x, y0-y, c)
            img.Set(x0-x, y0-y, c)
            img.Set(x0+y, y0+x, c)
            img.Set(x0-y, y0+x, c)
            img.Set(x0+y, y0-x, c)
            img.Set(x0-y, y0-x, c)
            if f >= 0 {
                dy += 2
                f += dy
            }
            dx += 2
            f += dx
            x++
        }
    }
    func parseHexColor(s string) (color.RGBA, error) {
        // ... implementation ...
    }