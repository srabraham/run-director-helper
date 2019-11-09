package parkrun

import (
	"strings"
	"testing"
)

// Basic happy path test
func TestLastEventNumberSuccess(t *testing.T) {
	html := `
<div class="Results" data-name="Name" data-agegroup="Age Group" data-club="Club" data-gender="Gender" data-achievement="Achievement" data-unknown="Unknown">
	<div class="Results-header">
		<h1>South Boulder Creek parkrun</h1>
		<h3>09/11/2019<span class="spacer"> | </span><span>#89</span></h3>
	</div>
</div>
`

	lastEventNum, err := lastEventNumber(strings.NewReader(html))
	if err != nil {
		t.Error(err)
	}
	if lastEventNum != 89 {
		t.Errorf("Expected lastEventNumber 89, instead %v", lastEventNum)
	}
}
