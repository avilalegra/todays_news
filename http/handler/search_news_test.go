package handler

import (
	"avilego.me/recent_news/news"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewPreviewData(t *testing.T) {
	preview := previews[0]
	prvData := newPreviewData(preview)
	assert.Equal(t, previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}, prvData)
}

func TestNewSearchResponse(t *testing.T) {
	for _, tData := range tsMakeSearchResponse {
		response := newSearchResponse(tData.previews)
		assert.Equal(t, tData.response, response)
	}
}

var tsMakeSearchResponse = []struct {
	previews []news.Preview
	response searchResponse
}{
	{
		previews[:2],
		searchResponse{
			Count: 2,
			Data: searchData{
				Sources: []news.Source{
					*sources["phoronix"],
				},
				Previews: []previewData{
					newPreviewData(previews[0]),
					newPreviewData(previews[1]),
				},
			},
		},
	},
	{
		previews[:3],
		searchResponse{
			Count: 3,
			Data: searchData{
				Sources: []news.Source{
					*sources["phoronix"],
					*sources["rtve"],
				},
				Previews: []previewData{
					newPreviewData(previews[0]),
					newPreviewData(previews[1]),
					newPreviewData(previews[2]),
				},
			},
		},
	},
	{
		previews[0:0],
		searchResponse{
			Count: 0,
			Data: searchData{
				Sources:  []news.Source{},
				Previews: []previewData{},
			},
		},
	},
}

type finderMock struct {
	t        *testing.T
	keywords string
	previews []news.Preview
}

func (b finderMock) Find(keywords string) []news.Preview {
	if b.keywords != keywords {
		assert.FailNow(b.t, fmt.Sprintf("finder called with keywords '%v', but '%v' was expected", keywords, b.keywords))
	}
	return b.previews
}

func TestSearch(t *testing.T) {
	for _, tData := range tsSearch {
		handler := searchHandler{finderMock{t, tData.keywords, tData.previews}}
		params := url.Values{}
		params.Set("keywords", tData.keywords)

		expectedJson, _ := json.Marshal(newSearchResponse(tData.previews))

		assert.HTTPBodyContains(t, handler.ServeHTTP, "GET", "/search", params, string(expectedJson))
	}
}

var tsSearch = []struct {
	keywords string
	previews []news.Preview
}{
	{
		"AMD",
		previews[0:2],
	},
	{
		"amd",
		previews[0:2],
	},
	{
		"",
		previews[0:0],
	},
}

var sources = map[string]*news.Source{
	"phoronix": {
		Title:       `Phoronix`,
		Link:        `https://www.phoronix.com/`,
		Language:    `en-US`,
		Description: `Linux Hardware Reviews & News`,
	},
	"rtve": {
		Title:       `Noticias en rtve.es`,
		Link:        `http://www.rtve.es`,
		Description: `RSS Tags`,
	},
}

var previews = []news.Preview{
	{
		Title:       `AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01`,
		Description: `Since July we've seen AMD open-source driver engineers posting code for "Cyan Skillfish" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...`,
		Source:      sources["phoronix"],
	},
	{
		Title:       `Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
		Description: `A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...`,
		Source:      sources["phoronix"],
	},
	{
		Title:       `Erupción en La Palma, en directo | La lava llega a 800 metros del mar y cambia de dirección al norte`,
		Link:        `http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml`,
		Description: `<ul> <li>Varios n&uacute;cleos poblacionales del municipio de Tazacorte han sido confinados</li> <li>La colada de lava podr&iacute;a llegar a la costa en las pr&oacute;ximas horas</li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TES800089&si=http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml" alt=""/>`,
		Source:      sources["rtve"],
	},
	{
		Title:       `Guía de restricciones COVID: nuevas medidas en ocio nocturno, hostelería y aforos, directo`,
		Link:        `http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml`,
		Description: `<ul> <li>Repasa las principales medidas y restricciones frente a la COVID-19, comunidad a comunidad del municipio</li> <li><a href="https://www.rtve.es/noticias/20210928/coronavirus-covid-directo-espana-mundo-ultima-hora/2175601.shtml" target="_blank">Coronavirus: &uacute;ltima hora</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-del-coronavirus-espana/2004681.shtml" target="_blank">Mapa de Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/ocupacion-camas-covid-19-hospitales-espanoles/2042349.shtml" target="_blank">Hospitales y UCI</a></li> <li><a href="https://www.rtve.es/noticias/20210924/campana-vacunacion-espana/2062499.shtml" target="_blank">Vacunas en Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-mundial-del-coronavirus/1998143.shtml" target="_blank">Mapa mundial&#8203;</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/lab/vacunacion-espana-coronavirus/">Especial: La gran vacunaci&oacute;n</a></li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TELCO20VX&si=http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml" alt=""/>`,
		Source:      sources["rtve"],
	},
}