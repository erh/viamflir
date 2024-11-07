package viamflir

import (
	"testing"

	"go.viam.com/test"
)

var exampleXML = `<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
  <specVersion>
    <major>1</major>
    <minor>0</minor>
  </specVersion>
  <device>
    <deviceType>urn:schemas-upnp-org:device:basic:1</deviceType>
    <friendlyName>M364C TAVY16D</friendlyName>
    <manufacturer>FLIR Systems, Inc.</manufacturer>
    <manufacturerURL>http://www.flir.com</manufacturerURL>
    <modelDescription>FLIR Infrared Camera</modelDescription>
    <modelName>M364C</modelName>
    <modelNumber>E70518</modelNumber>
    <modelURL>http://www.flir.com/marine/</modelURL>
    <serialNumber>TAVY16D</serialNumber>
    <UDN>uuid:IRCamera-1_0-89911a22-ef16-11dd-84a7-0011c7156eab</UDN>
    <UPC></UPC>
    <serviceList>
    </serviceList>
   <presentationURL>http://172.16.5.12:80/</presentationURL>
   <buildDate>00/00/00</buildDate>
   <deviceControl>JCU</deviceControl>
   <submodelName></submodelName>
   <deviceInfo>
     <CGIport>8090</CGIport>
     <ONVIFport>8091</ONVIFport>
     <productId></productId>
   </deviceInfo>
</device>
</root>`

func TestParse1(t *testing.T) {

	dd, err := parseDeciceDesc("", []byte(exampleXML))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, dd.SpecVersion.Major, test.ShouldEqual, 1)
	test.That(t, dd.Device.ModelName, test.ShouldEqual, "M364C")
	test.That(t, dd.Device.Manufacturer, test.ShouldEqual, "FLIR Systems, Inc.")
	test.That(t, dd.isFLIR(), test.ShouldBeTrue)
}

func TestPath1(t *testing.T) {
	test.That(t, Config{}.path(), test.ShouldEqual, "/vis.1")
	test.That(t, Config{"/vis.1"}.path(), test.ShouldEqual, "/vis.1")
	test.That(t, Config{"vis.1"}.path(), test.ShouldEqual, "/vis.1")
}
