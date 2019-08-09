package tests

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"strconv"
	"time"
)

const (
	DATE_FORMAT        = "2006-01-02 15:04:05Z"
	TEST_ENDPOINT      = "http://localhost:9010/v1/event"
	CONTENT_TYPE       = "Content-Type"
	CONTENT_TYPE_CE    = "application/cloudevents"
	PUBLISH_DATE       = "publishDate"
	PUBLISH_QUEUE      = "publishQueue"
	PUBLISH_TYPE       = "publishtype"
	TEST_PUBLISH_QUEUE = "throwAt"
	TEST_PUBLISH_TYPE  = "data_only"
)

func CreateEventRequest() {
	ginkgo.It("Valid msg", func() {
		for i := 1; i < 300; i++ {
			strIvalue := strconv.Itoa(i)
			fmt.Print("Trying to create an event " + strIvalue)
			mockReader, err := GetMockReader(getMockEvent(time.Now(), strIvalue))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			h := make(map[string]string)
			h[CONTENT_TYPE] = CONTENT_TYPE_CE
			h[PUBLISH_DATE] = time.Now().Add(time.Duration(i) * 100 * time.Millisecond).UTC().Format(DATE_FORMAT)
			h[PUBLISH_QUEUE] = TEST_PUBLISH_QUEUE
			h[PUBLISH_TYPE] = TEST_PUBLISH_TYPE
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		}
	})
}

func getMockEvent(publihsDate time.Time, ref string) interface{} {
	date, err := json.Marshal(struct {
		SpecVersion  string `json:"specversion"`
		Type         string `json:"type"`
		Source       string `json:"source"`
		ID           string `json:"id"`
		Time         string `json:"time,omitempty"`
		PublishDate  string `json:"publishdate"`
		PublishQueue string `json:"publishqueue"`
		PublishType  string `json:"publishtype"`
		Data         string `json:"data"`
	}{
		//"source": "sourceEvent",
		//"id": "78b818ae-7d353-11de9-u98dc1-54f64fr71frdg2d",
		//"time": "2019-05-23T12:08:27.250749658Z",
		//"publishDate": "2019-07-01 17:37:00Z",
		//"publishQueue": "throwAt",
		//"data": {"errorQueue":"Error.Timer.Resource.Scheduled","retry":3,"type":"timer","referenceName":"Timer.Resource.ThrowAt","message":{"referenceName":"Timer.Resource.ThrowAt","origin":[],"data":{"eventDate":{"date":"2019-05-13 19:20:04.000000","timezone_type":3,"timezone":"UTC"},"event":{"referenceName":"Native.Location.Expire","origin":[],"data":{"eventDate":{"date":"2019-05-13 19:20:13.000000","timezone_type":3,"timezone":"UTC"},"location":{"latitude":-24.812354,"longitude":-13.251654,"velocity":null,"precision":null,"direction":null,"address":null,"creationDate":{"date":"2019-05-13 19:20:03.000000","timezone_type":3,"timezone":"UTC"},"locationDate":{"date":"2019-05-13 19:20:08.000000","timezone_type":3,"timezone":"UTC"},"timeLapse":null,"precisionQualifier":null,"reason":null,"source":"RadioGateway","rejected":false,"rejectedReason":null,"driver":null,"odometer":null,"horimetre":null,"batteryLevel":null},"user":{"userId":1,"name":"TestUser 1","active":true,"userType":{"userTypeId":1,"name":"TestUserType 1","statusSet":{"statusSetId":1,"name":"Set 1","defaultDescription":"Set 1","defaultMapImage":"\/file\/fcef71cb","inactiveMapImage":"\/file\/f2a5ce77","unknownMapImage":"\/file\/331f50e4","unknownTimeToLive":1,"hasUnknownAlarm":false,"hasInactiveAlarm":false,"unknownAlarmMessage":"","inactiveAlarmMessage":"","generateUnknownStatus":true,"statusOverflow":null,"statusSetItemCollection":[{"statusSetItemId":1,"code":32768,"description":"Status 0","image":"\/file\/3be52bd8","timeToLive":null,"hasAlarm":false,"alarmMessage":null,"statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":2,"code":32769,"description":"Status 1","image":"\/file\/3c246e41","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":1,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":3,"code":32770,"description":"Status 2","image":"\/file\/3e32ef25","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":2,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":4,"code":32771,"description":"Status 3","image":"\/file\/3fa4b22c","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":3,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":5,"code":32772,"description":"Status 4","image":"\/file\/4aacc9df","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":4,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":6,"code":32773,"description":"Status 5","image":"\/file\/4af0a333","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":5,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":7,"code":32774,"description":"Status 6","image":"\/file\/4c6e8bdb","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":8,"code":32775,"description":"Status 7","image":"\/file\/4d6d0de8","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":9,"code":32776,"description":"Status 8","image":"\/file\/4f82ae0c","timeToLive":1,"hasAlarm":false,"alarmMessage":null,"statusOverflow":8,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":20,"code":32777,"description":"Status 9","image":"\/file\/8c28de6c","timeToLive":null,"hasAlarm":true,"alarmMessage":"5000","statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":21,"code":32778,"description":"Status 10","image":"\/file\/8c35deeb","timeToLive":null,"hasAlarm":false,"alarmMessage":"5000","statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":22,"code":32779,"description":"Status 11","image":"\/file\/10c1f975","timeToLive":null,"hasAlarm":true,"alarmMessage":"","statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":23,"code":32780,"description":"Status 12","image":"\/file\/13efb338","timeToLive":null,"hasAlarm":true,"alarmMessage":"g76y7u","statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":24,"code":32781,"description":"Status 13","image":"\/file\/28b65af2","timeToLive":null,"hasAlarm":false,"alarmMessage":"","statusOverflow":null,"generateSameStatus":true,"hasEvent":false},{"statusSetItemId":25,"code":32782,"description":"Status 14","image":"\/file\/28deb28b","timeToLive":null,"hasAlarm":false,"alarmMessage":"g76y7u","statusOverflow":null,"generateSameStatus":true,"hasEvent":true}]}},"observation":""},"device":{"deviceId":"5000","deviceUniqueId":8,"system":{"systemId":2,"description":"SystemTest","code":"724-1121","systemType":{"systemTypeId":2,"description":"SystemTypeTest","code":"Tetra"}}},"contract":{"contractId":1,"contractCode":"0000\/0000","contractEmail":"test@email.com","enterprise":{"enterpriseId":1,"enterpriseDescription":"Tester Enterprise","enterpriseCode":"000000","enterpriseIdentification":"00000000000000"}},"expirationDate":{"date":"2019-05-13 19:20:13.000000","timezone_type":3,"timezone":"UTC"},"version":3},"context":[]},"reference":"Native.Location.Expire","throwAt":1557775213,"throwInMilliseconds":1000,"iterations":8,"maxTimeWait":1000},"context":[]}}
		//
		//
		//
		//
		//publishDate: publihsDate.Format("teste"),
		//text:        fmt.Sprintf("Test event %s", ref),
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return date
}
