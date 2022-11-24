# Package Analyzer

## Solution for providing basic comparison function for ALT Linux repositories
PA connects to official [ALT Linux repositories storage](https://rdb.altlinux.org/api/export/branch_binary_packages/{branch}).
In current scope only branches "sisyphus" and "p10" are included. Comparison conducted for each branch (hereby source) against another one (hereby target).
PA compares 2 branches and:
* finds unique package in source, which are missing in target
* finds package present in both branches, but with higher releases in the source

Results are provided ass array of json strings in the format:
{<branch> <method> <architecture> <packages_count>}
- branch - sisyphus or p10
- method
    * "unique" means present in the branch only
    * "higher" - release in branch is higher than in the other one
    * "failure" - failures in comparing package info
* multile"unique" means present in the branch only, "higher" - release in branch is higher than in the other one
- architecture correspondes to the "arch" field in the package
- packages_count for number of packages for comparison method and architecture found in branch

__Sample output__
```
{sisyphus unique armh 2129}
...
{sisyphus unique x86_64-i586 783}
{sisyphus unique noarch 4431}
{p10 unique i586 2173}
...
{p10 unique noarch 2512}
{p10 unique x86_64-i586 908}
{sisyphus higher noarch 4336}
...
{sisyphus higher x86_64 11638}
{sisyphus higher x86_64-i586 4290}
{sisyphus failure  117}
{p10 higher armh 2018}
{p10 higher i586 2088}
...
{p10 higher x86_64-i586 756}
{p10 higher noarch 838}
{p10 failure  471}
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



