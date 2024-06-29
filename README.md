### prerequisites:
- golang 1.22
- [optional] ansible
- [optional] docker


### build:
- run `ansible-playbook tasks.yml`
- check `./build/`:
    - file `importer`: native executable
    - file `importer.zip`: zip of extensions, ready to upload in extension hub
    - run `docker images | importer`: a docker image of the executable
