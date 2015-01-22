package aqua

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {

	Convey("Given a RestServer", t, func() {
		s := NewRestServer()

		Convey("Then its embedded Fixture should have proper default values", func() {
			So(s.Url, ShouldBeEmpty)
			So(s.Root, ShouldBeEmpty)
			So(s.Version, ShouldEqual, "*")
		})
	})
}

func TestGetShouldNotHonourPost(t *testing.T) {
	s := NewRestServer()
	port := getUniquePortForTestCase()
	s.RunWith(port, false)

	Convey("Given a RestServer", t, func() {
		Convey("A url set for Http GET should return 404 for POST", func() {
			postData := make(map[string]string)
			postData["abc"] = "def"
			url := fmt.Sprintf("http://localhost:%d/aqua/ping", port)
			code, _, _ := postUrl(url, postData, nil)
			So(code, ShouldEqual, 404)
		})
	})
}

type AnyService struct {
	RestService
	honourGet  GetApi `url:"a-url"`
	honourPost GetApi `url:"a-url"`
}

func (me *AnyService) HonourGet() string  { return "" }
func (me *AnyService) HonourPost() string { return "" }

func TestSameUrlWithSameHttpMethods(t *testing.T) {
	s := NewRestServer()

	Convey("Given a RestServer", t, func() {
		Convey("When loading services with same urls and same http methods", func() {
			Convey("Then the program should panic", func() {
				So(func() {
					s.AddService(&AnyService{})
				}, ShouldPanic)
			})
		})
	})
}

type AnyServiceA struct {
	RestService
	honourGet  GetApi  `url:"a-url"`
	honourPost PostApi `url:"a-url"`
}

func (me *AnyServiceA) HonourGet() string  { return "" }
func (me *AnyServiceA) HonourPost() string { return "" }

func TestSameUrlWithDifferentHttpMethods(t *testing.T) {
	s := NewRestServer()

	Convey("Given a RestServer", t, func() {
		Convey("When loading services with same urls but different http methods", func() {
			Convey("Then the program should NOT panic", func() {
				So(func() {
					s.AddService(&AnyServiceA{})
				}, ShouldNotPanic)
			})
		})
	})
}

func TestAddMethodValidations(t *testing.T) {

	Convey("Given a RestServer", t, func() {
		s := NewRestServer()

		type BasicService struct {
			RestService
		}
		svc := BasicService{}

		Convey("When a Service object is directly passed to its add method", func() {
			Convey("Then there should be panic", func() {
				So(func() {
					s.AddService(svc)
				}, ShouldPanic)
			})
		})

		Convey("When address of Service object is passed to its add method", func() {
			Convey("Then it gets well accepted", func() {
				So(func() {
					s.AddService(&svc)
				}, ShouldNotPanic)
			})
		})

		Convey("When address of an object is passed to it", func() {
			Convey("And the object is not composed of RestService struct", func() {
				type IsNotComposedOfRestService struct{}
				obj := IsNotComposedOfRestService{}
				Convey("Then there should be panic", func() {
					So(func() {
						s.AddService(&obj)
					}, ShouldPanic)
				})
			})
		})

		Convey("When address an object is passed to it", func() {
			Convey("And the object is composed of RestService struct in a named (non-anonymous) field", func() {
				type IsComposedOfNamedRestService struct {
					field RestService
				}
				obj := IsComposedOfNamedRestService{}
				Convey("Then there should be panic", func() {
					So(func() {
						s.AddService(&obj)
					}, ShouldPanic)
				})
			})
		})
	})
}

type UserServiceA struct {
	RestService `root:"/A"`
	getUser     GetApi
}

func (me *UserServiceA) GetUser() string { return "" }

type UserServiceB struct {
	RestService `version:"0.3" root:"/B"`
	getUser     GetApi
}

func (me *UserServiceB) GetUser() string { return "" }

type UserServiceC struct {
	RestService `version:"0.3" root:"/C"`
	getUser     GetApi
}

func (me *UserServiceC) GetUser() string { return "" }

type UserServiceD struct {
	RestService `version:"0.3" root:"/D"`
	getUser     GetApi `version:"0.5"`
}

func (me *UserServiceD) GetUser() string { return "" }

func TestConfigurationHierarchy(t *testing.T) {

	Convey("Given a RestServer", t, func() {
		s := NewRestServer()
		s.Version = "0.2"

		Convey("Then the server config is inherited at Fixture", func() {
			u := UserServiceA{}
			s.AddService(&u)
			apiId := getServiceId("GET", "", "0.2", "/A/get-user")
			_, found := s.apis[apiId]
			So(found, ShouldBeTrue)
		})

		Convey("Then the service tag overrides server config", func() {
			u := UserServiceB{}
			s.AddService(&u)
			apiId := getServiceId("GET", "", "0.3", "/B/get-user")
			_, found := s.apis[apiId]
			So(found, ShouldBeTrue)
		})

		Convey("Then programmatically set service values override service tag", func() {
			u := UserServiceC{}
			u.Version = "0.4"
			s.AddService(&u)
			apiId := getServiceId("GET", "", "0.4", "/C/get-user")
			_, found := s.apis[apiId]
			So(found, ShouldBeTrue)
		})

		Convey("Then the Fixture tag overrides service values", func() {
			u := UserServiceD{}
			u.Version = "0.4"
			s.AddService(&u)
			apiId := getServiceId("GET", "", "0.5", "/D/get-user")
			_, found := s.apis[apiId]
			So(found, ShouldBeTrue)
		})
	})
}
