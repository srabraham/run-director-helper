package parkrun

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

const (
	futureRosterContents = `
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

	// This is very incomplete. Fill in the gaps as needed.
	latestResultsContents = `
	</div></div><?xml version="1.0"?>
	<h2>
      South Boulder Creek parkrun #
			52 -
			23/02/2019</h2>
			<p><a href="#explainTable"></a></p>
			<p><table class="sortable" id="results" align="center">
			    <thead></thead>
			    </table></p>
			`
)

func StartTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Got request: %+v", r)
			log.Printf("Got request at URL: %s", r.URL)
			switch r.Method {
			case "GET":
				if r.URL.String() == "/futureroster/" {
					fmt.Fprint(w, futureRosterContents)
					return
				}
				if r.URL.String() == "/results/latestresults/" {
					fmt.Fprintf(w, latestResultsContents)
					return
				}
			}
			w.WriteHeader(http.StatusNotImplemented)
		}))
}
