package parkrun

import (
	"strings"
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	html := `
<h2>Example text!</h2>
<table id="rosterTable">
   <thead>
      <tr>
         <th> </th>
         <th>2 March 2019</th>
         <th>9 March 2019</th>
         <th>16 March 2019</th>
         <th>23 March 2019</th>
         <th>30 March 2019</th>
         <th>6 April 2019</th>
      </tr>
   </thead>
   <tbody>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Run Director</a></th>
         <td>Peter PEPPER</td>
         <td>Peter PEPPER</td>
         <td>Sally SIMMONS</td>
         <td>Sally SIMMONS</td>
         <td>Tommy THOMAS</td>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Equipment Storage and Delivery</a></th>
         <td>Peter PEPPER</td>
         <td>Peter PEPPER</td>
         <td/>
         <td>Sally SIMMONS</td>
         <td/>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Timekeeper</a></th>
         <td>Rod ROMAN</td>
         <td/>
         <td/>
         <td/>
         <td/>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Barcode Scanning</a></th>
         <td>Peter PEPPER</td>
         <td>Peter PEPPER</td>
         <td>Sally SIMMONS</td>
         <td>Sally SIMMONS</td>
         <td>Tommy THOMAS</td>
         <td>Gary GERRIT</td>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Finish Tokens</a></th>
         <td>Matty MATTHEWS</td>
         <td/>
         <td/>
         <td/>
         <td/>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Tail Walker</a></th>
         <td>Sally SIMMONS</td>
         <td>Timmy TIMMONS</td>
         <td>Greg GRATH</td>
         <td>Greg GRATH</td>
         <td/>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Marshal</a></th>
         <td>Nick NORMAN</td>
         <td>Greg GRATH</td>
         <td/>
         <td/>
         <td/>
         <td/>
      </tr>
      <tr>
         <th bgcolor="#FFFFAA"><a href="http://example.com" class="voltask" target="_blank" title="click for description">Marshal</a></th>
         <td>Lizzy LYONS</td>
         <td>Donny DONALDS</td>
         <td/>
         <td/>
         <td/>
         <td/>
      </tr>
      <tr>
         <th><a href="http://example.com" class="voltask" target="_blank" title="click for description">Results Processor</a></th>
         <td>Peter PEPPER</td>
         <td>Peter PEPPER</td>
         <td>Sally SIMMONS</td>
         <td>Sally SIMMONS</td>
         <td>Tommy THOMAS</td>
         <td/>
      </tr>
   </tbody>
</table>`

	roster, err := fetchFutureRoster(strings.NewReader(html))
	if err != nil {
		t.Error(err)
	}
	loc, err := time.LoadLocation(*eventLocation)
	if err != nil {
		t.Error(err)
	}
	res1 := roster.SortedEvents[1]
	expectedDateMidnight := time.Date(2019, time.March, 9, 0, 0, 0, 0, loc)
	expectedDate := expectedDateMidnight.Add(*eventTime)
	if !res1.Date.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, res1.Date)
	}
	if res1.RoleVolunteers[0].Role != "Run Director" {
		t.Errorf("Wrong! %v", res1.RoleVolunteers[0])
	}
	if res1.RoleVolunteers[0].Volunteer != "Peter PEPPER" {
		t.Errorf("Wrong! %v", res1.RoleVolunteers[0])
	}
	if res1.RoleVolunteers[2].Role != "Timekeeper" {
		t.Errorf("Wrong! %v", res1.RoleVolunteers[2])
	}
	if res1.RoleVolunteers[2].Volunteer != "" {
		t.Errorf("Wrong! %v", res1.RoleVolunteers[2])
	}
}

func assertNextEventTime(t *testing.T, fr FutureRoster, afterTime time.Time, expectedTime time.Time) {
	testEvent, err := fr.FirstEventAfter(afterTime)
	if err != nil {
		t.Error(err)
	}
	if testEvent.Date != expectedTime {
		t.Errorf("Expected testEvent.Date == %v, got %v", expectedTime, testEvent.Date)
	}
}

func TestFirstEventAfter(t *testing.T) {
	loc, err := time.LoadLocation("America/Inuvik")
	if err != nil {
		t.Error(err)
	}
	time0 := time.Date(2012, time.November, 12, 9, 0, 0, 0, loc)
	time1 := time.Date(2012, time.November, 19, 9, 0, 0, 0, loc)
	fr := FutureRoster{[]EventDetails{
		{Date: time0},
		{Date: time1},
	}}

	beforeTime0 := time0.Add(time.Second * -1)
	assertNextEventTime(t, fr, beforeTime0, time0)

	isTime0 := time0
	assertNextEventTime(t, fr, isTime0, time1)

	afterTime0 := time0.Add(time.Second * 1)
	assertNextEventTime(t, fr, afterTime0, time1)
}
