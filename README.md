# Package Analyzer

## Solution for providing basic comparison function for ALT Linux repositories
PA connects to official [ALT Linux repositories storage](https://rdb.altlinux.org/api/export/branch_binary_packages/{branch}).
In current scope only branches "sisyphus" and "p10" are included. Comparison conducted for each branch (hereby source) against another one (hereby target).
PA compares 2 branches and:
* finds unique package in source, which are missing in target
* finds package present in both branches, but with higher releases in the source

Results are provided within *.json files named by architecture.
JSON structure for each json file is as follows:
```
type ResultsOutput3 struct {
  Method []struct {
     Branch []struct {
      Packages []Binpack `json:"packages"`
     } `json:"branch"`
  } `json:"method"`
}
```

---

### Building instructions
`cd cmd/agent`
`go build -o <execuitable>`

### Running application
./<executable> [-s <scope>]

-s is optional flag for limiting scope of data analysis. Possible options:
* all - finds unique packages and higher releases in branches
* diff - only finds unique packages in branches
* releases - only finds higher releases in branches



