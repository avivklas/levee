package levee

import (
	"testing"
	"time"
)

type takeTestRequest struct {
	timeFromStart	time.Duration
	amount			int64
	expectedResult	time.Duration
}

var takeTests = []struct{
	description	string
	rate		float64
	requests	[]takeTestRequest
} {{
	description:	"1/1sec",
	rate:			1,
	requests: 		[]takeTestRequest{{
		timeFromStart:	0,
		amount:			0,
		expectedResult: 0,
	},{
		timeFromStart:	0,
		amount:			1,
		expectedResult: 1 * time.Second,
	},{
		timeFromStart:	2 * time.Second,
		amount:			1,
		expectedResult: 0,
	},{
		timeFromStart:	30 * time.Second,
		amount:			30,
		expectedResult: 2 * time.Second,
	}},
},{
	description:	"1/2sec",
	rate:			0.5,
	requests: 		[]takeTestRequest{{
		timeFromStart:	0,
		amount:			0,
		expectedResult: 0,
	},{
		timeFromStart:	0,
		amount:			1,
		expectedResult: 2 * time.Second,
	},{
		timeFromStart:	2 * time.Second,
		amount:			1,
		expectedResult: 2 * time.Second,
	},{
		timeFromStart:	30 * time.Second,
		amount:			15,
		expectedResult: 2 * 2 * time.Second,
	}},
},{
	description:	"2/1sec",
	rate:			2,
	requests: 		[]takeTestRequest{{
		timeFromStart:	0,
		amount:			0,
		expectedResult: 0,
	},{
		timeFromStart:	0,
		amount:			1,
		expectedResult: 500 * time.Millisecond,
	},{
		timeFromStart:	1 * time.Second,
		amount:			1,
		expectedResult: 0,
	},{
		timeFromStart:	30 * time.Second,
		amount:			60,
		expectedResult: 2 * 500 * time.Millisecond,
	}},
},{
	description:	"8/1sec",
	rate:			8,
	requests: 		[]takeTestRequest{{
		timeFromStart:	0,
		amount:			0,
		expectedResult: 0,
	},{
		timeFromStart:	0,
		amount:			1,
		expectedResult: 125 * time.Millisecond,
	},{
		timeFromStart:	1 * time.Second,
		amount:			8,
		expectedResult: 125 * time.Millisecond,
	},{
		timeFromStart:	30 * time.Second,
		amount:			240,
		expectedResult: 9 * 125 * time.Millisecond,
	}},
},{
	description:	"1/8sec",
	rate:			0.125,
	requests: 		[]takeTestRequest{{
		timeFromStart:	0,
		amount:			0,
		expectedResult: 0,
	},{
		timeFromStart:	0,
		amount:			1,
		expectedResult: 8 * time.Second,
	},{
		timeFromStart:	1 * time.Second,
		amount:			1,
		expectedResult: 15 * time.Second,
	},{
		timeFromStart:	30 * time.Second,
		amount:			4,
		expectedResult: 18 * time.Second,
	}},
},{
	description:	"deviation",
	rate:			0.125,
	requests: 		[]takeTestRequest{{
		timeFromStart:	7 * time.Second + 600 * time.Millisecond,
		amount:			1,
		expectedResult: 0 * time.Millisecond,
	},{
		timeFromStart:	15 * time.Second + 400 * time.Millisecond,
		amount:			1,
		expectedResult: 600 * time.Millisecond,
	}},
}}

func TestTake(t *testing.T) {
	for _, test := range takeTests {
		b := NewBucket(test.rate, 1000, time.Unix(0, 0))
		for j, req := range test.requests {
			d := b.take(b.firstTickTime.Add(req.timeFromStart), req.amount)
			if d != req.expectedResult {
				t.Fatalf("[%s test#%d] expected: %v actual: %v", test.description, j, req.expectedResult, d)
			}
		}
	}
}

