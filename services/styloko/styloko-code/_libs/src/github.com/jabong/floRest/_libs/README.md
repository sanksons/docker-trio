Contains the external libs.

Developer should manually go get <their_external_dependency>
And copy the folder inside _libs/src/

Till the time we find a better dependency management,
the external libs should be checked in discovery repository.

Delete the .git/ and .gitignore

The list of external dependencies:
```
github.com/BurntSushi/toml	
github.com/bradfitz/gomemcache	
github.com/natefinch/lumberjack	
github.com/onsi	
github.com/yvasiyarov
gopkg.in/yaml.v1
bitbucket.org/inflect
```
