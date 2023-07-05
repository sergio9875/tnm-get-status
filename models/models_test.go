package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestEvent_MarshalJSON(t *testing.T) {
	event := &Event{
		Action:      "Test",
		RequestedBy: 123,
		RequestedTS: time.Date(2021, 1, 26, 1, 1, 1, 0, time.Local),
	}
	data, err := json.Marshal(event)
	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
		return
	}
	if string(data) != "{\"requested_at\":1611615661,\"action\":\"Test\",\"requested_by\":123}" {
		t.Errorf("Error while marshalling: RequestedBy expected[%s], got[%s]",
			"{\"requested_at\":1611615661,\"action\":\"Test\",\"requested_by\":123}",
			string(data))
	}
}

func TestEvent_UnmarshalJSON(t *testing.T) {
	data := []byte("{\"requested_at\":1611615661,\"action\":\"Test\",\"requested_by\":123}")
	event := &Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		t.Errorf("Error while unmarshalling: %v", err)
		return
	}
	if event.Action != "Test" {
		t.Errorf("Error while unmarshalling: Action expected[%s], got[%s]", "Test", event.Action)
	}
	if event.RequestedBy != 123 {
		t.Errorf("Error while unmarshalling: RequestedBy expected[%d], got[%d]", 123, event.RequestedBy)
	}
	if event.RequestedTS.Unix() != 1611615661 {
		t.Errorf("Error while unmarshalling: RequestedAt expected[%d], got[%d]", 1611615661,
			event.RequestedTS.Unix())
	}
}

func TestEvent_UnmarshalJSONError(t *testing.T) {
	data := []byte("{\"requested_by\":\"azsx\"}")
	event := &Event{}
	err := json.Unmarshal(data, &event)
	if err == nil {
		t.Errorf("Error was expecting error while unmarshalling, got nil")
		return
	}
}

func TestSecretModel_Merge(t *testing.T) {
	sm := &SecretModel{}
	ism := aws.String(`{
    "db": {
      "treasury": {
        "dialect": "mysql"
      },,,,,,,,,,,,,,,,,,,,,,,,
    }
  }`)
	// bad json
	sm = sm.Merge(ism)
	if sm != nil {
		t.Errorf("Error was expecting a nil")
		return
	}
	// =============================
	sm = &SecretModel{}
	ism = aws.String(`{
    "db": {
      "treasury": {
        "dialect": "mysql"
      }
    }
  }`)
	// good json
	sm = sm.Merge(ism)
	if sm == nil {
		t.Errorf("Error was expecting a not nil")
		return
	}
	// =============================
	sm = &SecretModel{
		Database: &DBConfig{},
	}
	ism = aws.String(`{
    "db": {
      "treasury": {
        "dialect": "mysql"
      }
    }
  }`)
	// good json
	sm = sm.Merge(ism)
	if sm == nil {
		t.Errorf("Error was expecting a not nil")
		return
	}
	// =============================
	sm = &SecretModel{
		Database: &DBConfig{
			Africainv: &MssqlConfig{},
		},
	}
	ism = aws.String(`{
    "db": {
      "treasury": {
        "dialect": "mysql"
      }
    }
  }`)
	// good json
	sm = sm.Merge(ism)
	if sm == nil {
		t.Errorf("Error was expecting a not nil")
		return
	}
	// =============================
	sm = &SecretModel{
		Database: &DBConfig{
			Africainv: &MssqlConfig{},
		},
	}
	ism = aws.String(`{
    "db": {
      "treasury": {
        "dialect": "mysql"
      }
    }
  }`)
	// good json
	sm = sm.Merge(ism)
	if sm == nil {
		t.Errorf("Error was expecting a not nil")
		return
	}
}
