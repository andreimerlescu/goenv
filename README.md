# Go ENV

This utility is designed to allow you to let you interact with `.env` files in a manner consistent with `.env.local`, 
`.env.develop`, `.env.production`, etc., and we have shorter hand notations for interacting with `.env` files. 

## Installation

```bash
go get github.com/andreimerlescu/goenv@latest
```

### Package `env`

The basis behind `env` is the `env` package, which can be used in other Go applications and is licensed as Apache 2.0.

```bash
go get -u github.com/andreimerlescu/goenv/env
```

The [README.md](/env/README.md) has more information about using the package.

## Usage

I wrote this package because sometimes I just need to work with an env file, even if it doesn't yet exist. This is 
particularly useful for input validation of env files and the ability to transform the syntax of the .env file into
multiple formats. The [test.sh](test.sh) file has the summary of arguments that are best to run for the package. 


```sh
go install github.com:andreimerlescu/goenv@latest
goenv -has -env HOSTNAME || goenv -write -env HOSTNAME -value "$(hostname)" || { echo ERROR && exit 1; }
```

## Testing

```log
andrei@GitHub:~/repos/goenv|master⚡ ⇒  make all
Summary generated: summaries/summary.2025.08.09.19.28.57.UTC.md
Building goenv binary...
Clean successful: ./bin/goenv-darwin-arm64
Build successful: ./bin/goenv-darwin-arm64

andrei@goenv.git:. ⚡ Test #1 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -has -env HOSTNAME
andrei@goenv.git:. ⚡ Test #2 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -has -env NON_EXISTENT
andrei@goenv.git:. ⚡ Test #3 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -is -env DATABASE -value test_data
andrei@goenv.git:. ⚡ Test #4 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -is -env DATABASE -value wrong_data
andrei@goenv.git:. ⚡ Test #5 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -print
AWS_REGION=us-west-2
OUTPUT=json
HOSTNAME=localhost
DBUSER=readonly
DBPASS=readonly
DATABASE=test_data
andrei@goenv.git:. ⚡ Test #6 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -json
{
  "AWS_REGION": "us-west-2",
  "DATABASE": "test_data",
  "DBPASS": "readonly",
  "DBUSER": "readonly",
  "HOSTNAME": "localhost",
  "OUTPUT": "json"
}
andrei@goenv.git:. ⚡ Test #7 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -yaml
---
AWS_REGION: "us-west-2" 
OUTPUT: "json" 
HOSTNAME: "localhost" 
DBUSER: "readonly" 
DBPASS: "readonly" 
DATABASE: "test_data" 
andrei@goenv.git:. ⚡ Test #8 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -toml
DBPASS: "readonly" 
DATABASE: "test_data" 
AWS_REGION: "us-west-2" 
OUTPUT: "json" 
HOSTNAME: "localhost" 
DBUSER: "readonly" 
andrei@goenv.git:. ⚡ Test #9 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -ini
[default]
OUTPUT = json
HOSTNAME = localhost
DBUSER = readonly
DBPASS = readonly
DATABASE = test_data
AWS_REGION = us-west-2
andrei@goenv.git:. ⚡ Test #10 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -xml
<?xml version="1.0" encoding="UTF-8"?>
<env>
   <AWS_REGION>us-west-2</AWS_REGION>
   <OUTPUT>json</OUTPUT>
   <HOSTNAME>localhost</HOSTNAME>
   <DBUSER>readonly</DBUSER>
   <DBPASS>readonly</DBPASS>
   <DATABASE>test_data</DATABASE>
</env>
andrei@goenv.git:. ⚡ Test #11 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -write -add -env NEW_KEY -value 'a new value'
andrei@goenv.git:. ⚡ Test #12 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -has -env NEW_KEY
andrei@goenv.git:. ⚡ Test #13 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -write -add -env HOSTNAME -value 'another-host'
andrei@goenv.git:. ⚡ Test #14 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -is -env HOSTNAME -value localhost
andrei@goenv.git:. ⚡ Test #15 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -write -rm -env OUTPUT
andrei@goenv.git:. ⚡ Test #16 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -not -has -env OUTPUT
andrei@goenv.git:. ⚡ Test #17 ⇒  ./bin/goenv-darwin-arm64 -file new.env -add -env HELLO -value world -write
andrei@goenv.git:. ⚡ Test #18 ⇒  ./bin/goenv-darwin-arm64 -file new.env -is -env HELLO -value world
andrei@goenv.git:. ⚡ Test #19 ⇒  ./bin/goenv-darwin-arm64 -file sample.env -v
v0.0.2
andrei@goenv.git:. ⚡ Test #20 ⇒  ./bin/goenv-darwin-arm64 -file non_existent_file.env -add -env FOO -value bar || echo "Test success because we expected an error here."                                                                                                      
andrei@goenv.git:. ⚡ Test #21 ⇒  ./bin/goenv-darwin-arm64 -file non_existent_file.env -add -env FOO -value bar -write
All 21 tests PASS!
NEW: bin/goenv-darwin-amd64
NEW: bin/goenv-darwin-arm64
NEW: bin/goenv-linux-amd64
NEW: bin/goenv-darwin-arm64
NEW: bin/goenv.exe
NEW: /Users/andrei/go/bin/goenv
```