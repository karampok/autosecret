# Extended ZTP
Simple controllers that they add for pull secret/bmh-secrets for ZTP deployments.

## Description
Pull secret is copied from secret openshift-config/pull-secret.
BMH secret is copied from secret openshift-config/bmh-secret.

## Run me
```
# reference secret in place and then
oc apply -k github.com/karampok/autosecret/config/default
```

## Note
Draft implementation to be used for CI/CD, testing.

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
