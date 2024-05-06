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
## Building
```
go build  -o build/main main.go
```
## Testing
```
ginkgo -r -v
```
* Recursive and verbose

## To-Do 
* Bug fixes
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
    * Move model parsing to puml.go
    * Update puml.go to add the starting dependency to the title of the graph
    * Add puml tests
    * Add package versions to assist with comparisons between OS versions 
    * Move puml diagram output to separate block to our log messages so any information about the diagram generation doesn't get mixed into the diagram itself
* Add Makefile with goals
  * Build
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
  * Hide already installed packages
    * Useful if you only need to know the dependencies that still need to be installed on a system
    * Add command line flag
    * Add command to check if installed
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