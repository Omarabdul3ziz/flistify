# Flistify

Is a simple tool that can create flists from scratch, run it locally & storing it on the 0hub registry.

### Create flist

As simple as writing a Dockerfile and building image from it.

- Create an Zerofile
  ```Dockerfile
  FROM ubuntu:jammy
  KERNEL linux-modules-extra-5.15.0-25-generic
  RUN apt update && apt install openssh-server curl cloud-init
  ENV FOO=bar
  ENTRYPOINT /start.sh
  ```
- Build the flist
  ```
  flistify build <zerofile>
  ```

### Run flist locally

Simpliy it will `arch_chroot` to the flist directory and run `ENTRYPOINT` or command

- Run
  ```
  flistify run <flist> <command>
  ```

### Store on 0Hub

It will tar your flist and using your provided HUBTOKEN will push it to 0Hub

- Push
  ```
  export HUBTOKEN="ey..."
  flistify push <flist>
  ```
