# Debendency
A tool to output debian package dependencies

## Parameters
Output from `go run main.go -help` at this point in time:
```
Usage of /tmp/go-build1711827394/b001/exe/main:
  -d    output dependencies as a diagram
  -o string
        output directory to save installer files to
  -p string
        .deb package name to calculate dependencies for
  -s    output dependencies as salt code

```
## Example usage
* `./debendency -p jq` fetch dependencies
* `./debendency -p jq -d` fetch dependencies and create an output diagram of the flow of dependencies
* `./debendency -p jq -d -e` fetch dependencies and create an output diagram of the flow of dependencies, **excluding** dependencies already installed.
* 

## Development Setup
On Windows you'll need to install the Ginkgo CLI:
```
go install github.com/onsi/ginkgo/v2/ginkgo@v2.12.0
```
On Linux you can make use of [mise](https://mise.jdx.dev/) via the `.mise.toml` file to install local tooling

## Building
```
go build  -o build/debendency main.go
```
## Testing
```
ginkgo -r -v
```
* Recursive and verbose

## To-Do 
* Bug fixes
  * Update `packageModel.GetPackageFilename` to account for versions with colons becoming `%3a` in filenames
    * version `1:4.4.10-10ubuntu4` becomes -> `libcrypt1_1%3a4.4.10-10ubuntu4_amd64.deb`
  * ~~Support cases where `Pre-Depends` exists for `dpkger.ParseDependencies`~~
  * ~~Ignore case where `Depends` has a dependency (python:any)~~
  * Support case where file already has been downloaded and we don't create a model as we cannot parse the parameters from the download output
    * dpkg -I might be needed to list these details if the download doesn't happen - don't know what the file name is
    * apt download will grab the latest available version
    * Perhaps we can use apt list to build a list of the latest version we do have if the download fails?
* Create command line flags
  * Mandatory vs Optional, can we support these
* Test command line flags
* Add testing command line flags to only run if on linux - integration tests
* Add integration tests
  * Move linux command non-mocked tests into the integrationtests package 
* Logging
  * Add standard out logger
  * Add standard error logger
  * Add verbose flag to set logging level to DEBUG
* Add puml output
  * ~~Move model parsing to puml.go~~
  * Update puml.go to add the starting dependency to the title of the graph
  * ~~Add puml tests~~
  * Case where package has no dependencies, should we show just the package e.g. dos2unix?
  * New line support, how to escape it in a golang string but have `\n` print in standard out?
  * Support dependency version conditions on relationship  e.g. python(>2.7) dot label syntax `[label = ">= 2.7"]`
  * Update diagram to include OS version in title and filename?
  * ~~Add package versions to assist with comparisons between OS versions~~  
  * Move to declare package and version in single place referencing just the package name in the diagram
  * Move puml diagram output to separate block to our log messages so any information about the diagram generation doesn't get mixed into the diagram itself
  * Update tests to use Google's cmp library
  * Switch puml generation to use golang templating
* Add Makefile with goals
  * ~~Build~~
  * Test
  * Integration test
  * Linting?
  * Formatting
  * Coverage
  * Go generate for fakes
* Supplementary Features
  * apt-get update before first download
  * Download cache
    * Default location of cache `~/debendency/cache`
    * Flag to specify location for the cache
    * Flag to delete download cache after use?
    * Save model structs to a json file
  * Versionless puml output
  * Puml Diff
    * Versionless
    * Versioned
  * Hide already installed packages
    * Useful if you only need to know the dependencies that still need to be installed on a system
    * ~~Add command line flag~~
    * ~~Add command to check if installed~~
    * Add integration tests
    * Add unit tests
    * Add interface to call to ensure model has installed boolean set
  * Salt output
    * Salt output tests
    * Important to remember here that only **offline** installers require dependencies
    * Online installer is only required for the **first** dependency, apt-get does everything else
    * This can be modelled as:
      * Map root method
      * Map dependencies method
      * With an option to using the dependencies template if no root template is specified
  * Puppet output