# floRest
florest code repo

Go installation link
https://golang.org/doc/install

Use the download link.

##Fetch the project
```
cd <GOPROJECTSPATH>/
git clone https://github.com/jabong/floRest
```

##Setup the application log
```

sudo mkdir /var/log/florest/
chown <user_name who is running ./florest> /var/log
```

##Make (build, install, test)
```

cd <GOPROJECTSPATH>/florest/
make deploy ENV=dev
cd <GOPROJECTSPATH>/bin/
./florest
```

##Example url
```
http://localhost:8080/florest/v1/hello
```

##Steps to bootstrap from floRest
```
make newapp NEWAPP="<your_project_name>"
```

##Steps to pull changes from floRest
```
floRest is available as library only. So just sync your _libs/src/github.com/floRest with github.com/jabong/floRest master branch.

##Wiki link
https://wiki.jira.rocket-internet.de/display/INDFAS/floRest+Framework
```

##How to raise PR in floRest
```
Refer the how to contribute section in https://wiki.jira.rocket-internet.de/display/INDFAS/florest-FAQ
```

##Run coverage for local tests
make coverage

##Run coverage against external tests
Step1: make coverall

Step2: cd bin/

Step3: ./floRest.test -test.coverprofile coverage.cov

Step4: Now run your external tests against this build and it will keep on calculating coverage. Once all the
tests have run, just stop the program, coverage.cov file will be generated in bin/.

Step5: Run  below command to convert coverage in html format.

go tool cover -html=coverage.cov -o coverage.html
