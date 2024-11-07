package viamflir

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/koron/go-ssdp"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/utils"

	"github.com/viam-modules/viamrtsp"
)

var FlirModel = family.WithModel("flir")

type Config struct {
	Path string
}

func (cfg Config) path() string {
	if cfg.Path == "" {
		return "/vis.1"
	}

	if cfg.Path[0] == '/' {
		return cfg.Path
	}

	return "/" + cfg.Path
}

func (cfg *Config) Validate(path string) ([]string, error) {
	return nil, nil
}

func init() {
	resource.RegisterComponent(
		camera.API,
		FlirModel,
		resource.Registration[camera.Camera, *Config]{
			Constructor: newFlir,
		})
}

func newFlir(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (camera.Camera, error) {
	ip, err := FindIP(ctx, logger)
	if err != nil {
		return nil, err
	}

	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, err
	}

	cc := viamrtsp.Config{
		Address: fmt.Sprintf("rtsp://%s:8554/%s", ip, newConf.path()),
	}

	ccc := resource.Config{
		Name:             conf.Name,
		API:              conf.API,
		Model:            conf.Model,
		Frame:            conf.Frame,
		LogConfiguration: conf.LogConfiguration,
		Attributes: utils.AttributeMap{
			"rtsp_address": cc.Address,
		},

		ConvertedAttributes: &cc,
	}

	return viamrtsp.NewRTSPCamera(ctx, deps, ccc, logger)
}

func FindIP(ctx context.Context, logger logging.Logger) (string, error) {

	list, err := ssdp.Search(ssdp.All, 3, "")
	if err != nil {
		return "", err
	}

	for _, srv := range list {
		logger.Debugf("found service (%s) at %s", srv.Type, srv.Location)
		if !mightBeFlir(srv) {
			logger.Debugf("\t skipping")
			continue
		}

		logger.Infof("found possible good service (%s) at %s", srv.Type, srv.Location)
		desc, err := readDeciceDesc(srv.Location)
		if err != nil {
			logger.Warnf("cannot read description %v", err)
			continue
		}

		logger.Infof("got description %v", desc)
		if desc.isFLIR() {
			logger.Infof("found flir %s", srv.Location)

			url, err := url.Parse(srv.Location)
			if err != nil {
				return "", err
			}

			pcs := strings.Split(url.Host, ":")
			return pcs[0], nil
		}

	}

	return "", fmt.Errorf("no flir found")
}

func mightBeFlir(srv ssdp.Service) bool {
	if !strings.HasPrefix(srv.Type, "urn:schemas-upnp-org:device:basic") {
		return false
	}

	// not sure if this makes sense
	if !strings.HasSuffix(srv.Location, "devicedesc.xml") {
		return false
	}

	return true
}

type deviceDesc struct {
	XMLName     xml.Name `xml:"root"`
	SpecVersion struct {
		Major int `xml:"major"`
		Minor int `xml:"minor"`
	} `xml:"specVersion"`
	Device struct {
		Manufacturer string `xml:"manufacturer"`
		ModelName    string `xml:"modelName"`
	} `xml:"device"`
}

func (dd *deviceDesc) isFLIR() bool {
	return strings.HasPrefix(dd.Device.Manufacturer, "FLIR")
}

func readDeciceDesc(url string) (*deviceDesc, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't fetch xml(%s): %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http fetch (%s) not ok: %v", url, err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't ready body from (%s): %v", url, err)
	}

	return parseDeciceDesc(url, data)
}

func parseDeciceDesc(url string, data []byte) (*deviceDesc, error) {
	var desc deviceDesc
	err := xml.Unmarshal(data, &desc)
	if err != nil {
		return nil, fmt.Errorf("bad xml from (%s): %v", url, err)
	}

	return &desc, nil
}
