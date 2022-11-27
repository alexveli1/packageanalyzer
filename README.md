# Package Analyzer

## Solution for providing basic comparison function for ALT Linux repositories
PA connects to official [ALT Linux repositories storage](https://rdb.altlinux.org/api/export/branch_binary_packages/{branch}).
In current scope only branches "sisyphus" and "p10" are included. Comparison conducted for each branch (hereby source) against another one (hereby target).
PA compares 2 branches and:
* finds unique package in source, which are missing in target
* finds package present in both branches, but with higher releases in the source

Results are provided within *.json files named by ARCHITECTURE.
JSON structure for each *.json file is as follows:
```
{
  <method>:{
        <branch>:[{
                "name": "bat-debuginfo",
                "epoch": 0,
                "version": "0.22.1",
                "release": "alt1",
                "arch": "ppc64le",
                "disttag": "sisyphus+308009.100.1.1",
                "buildtime": 1665136148,
                "source": "bat"
               },
               {
                    ...
               }],
        <branch>:[{}, ... {}]
        },
  <method>:{
        <branch>:[{ ... }],
        <branch>:[{ ... }]
  }
}
```
`<method>`
* higher - displays packages with higher releases for subsequent branch
* unique - displays packages which are unique for architecture in branch

`<branch>`
* p10, sisyphus
---

### Building instructions
`cd cmd/agent`
`go build -o <execuitable>`

### Running application
`./<executable> [-s <scope>]`
-s is optional flag for limiting scope of data analysis. Possible options for the flag are:
* all - finds unique packages and higher releases in branches
* diff - only finds unique packages in branches
* releases - only finds higher releases in branches



