### limitations

- only FROM ubuntu archive
- RUN can't have &&
- RUN uses bash


### Future work

- [ ] `FROM` should read from the 0hub registry and should be able to use the other flists as base
- [ ] `flistify run` should isolate the running process by applying namespaces & cgroups
