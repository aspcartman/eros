package storage

import (
	"image"
	"bytes"
	"github.com/aspcartman/exceptions"
	"github.com/nfnt/resize"
	"github.com/fogleman/primitive/primitive"
	"strings"
	"runtime"
	"image/jpeg"
	"net/http"
	"io/ioutil"
	"github.com/aspcartman/eros/env"
)


type Photo struct {
	ID          uint
	UserID      string
	Source      string
	URL         string
	Original    []byte
	Small       []byte
	Placeholder []byte
}

func (tx *Tx) AllPhotos() <-chan Photo {
	rows := tx.Query(`select id, user_id, source, url, original, small, placeholder from photos`)
	ch := make(chan Photo, 100)
	go func(){
		defer close(ch)
		p := Photo{}
		for rows.Next() {
			rows.Scan(&p.ID, &p.UserID, &p.Source, &p.URL, &p.Original, &p.Small, &p.Placeholder)
			ch <- p
		}
 	}()
	return ch
}

func (tx *Tx) SavePhoto(p *Photo) {
	if p.ID == 0 {
		tx.Exec(`INSERT INTO photos (user_id, source, url, original, small, placeholder)
				 VALUES ($1,$2,$3,$4,$5,$6)`,
			p.UserID, p.Source, p.URL, p.Original, p.Small, p.Placeholder)
	} else {
		tx.Exec(`UPDATE photos
				 SET	(user_id, source, url, original, small, placeholder) = ($2,$3,$4,$5,$6)
				 WHERE 	id = $1`,
			p.ID, p.UserID, p.Source, p.URL, p.Original, p.Small, p.Placeholder)
	}

}


func (p *Photo) Refresh() {
	data := p.Original
	if len(data) == 0 {
		env.Log.WithField("url", p.URL).Info("downloading photo")
		data = get(p.URL)
	}

	env.Log.WithField("url", p.URL).Info("decoding photo")
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		e.Throw("failed decoding image", err, e.Map{
			"data": string(data),
		})
	}

	p.Original = data

	env.Log.WithField("url", p.URL).Info("making small")
	img, p.Small = small(img)

	env.Log.WithField("url", p.URL).Info("making placeholder")
	p.Placeholder = placeholder(img)
}

func get(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		e.Throw("failed getting url", err, e.Map{
			"url": url,
		})
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.Throw("failed reading response body", err, e.Map{
			"url": url,
		})
	}

	if res.StatusCode != 200 {
		e.Throw("bad status code", nil, e.Map{
			"url":  url,
			"code": res.StatusCode,
		})
	}

	return data
}

func small(img image.Image) (image.Image, []byte) {
	img = resize.Thumbnail(256, 256, img, resize.Bicubic)
	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, img, &jpeg.Options{80})
	if err != nil {
		e.Throw("failed encoding", err)
	}
	return img, buf.Bytes()
}

func placeholder(img image.Image) []byte {
	if img.Bounds().Size().X > 256 || img.Bounds().Size().Y > 256 {
		img = resize.Thumbnail(256, 256, img, resize.Bicubic)
	}

	bg := primitive.MakeColor(primitive.AverageImageColor(img))
	model := primitive.NewModel(img, bg, 256, runtime.NumCPU())
	for i := 0; i < 100; i++ {
		model.Step(primitive.ShapeTypeTriangle, 128, 0)
	}
	svg := model.SVG()
	svg = strings.Replace(svg, "<g", `<filter id="b"><feGaussianBlur stdDeviation="12" /></filter><g filter="url(#b)"`, 1)
	return []byte(svg)
}
