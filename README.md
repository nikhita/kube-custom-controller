# kube-custom-controller

This repository contains an example of how you can create a custom controller for custom resources in Kubernetes.

Whenever you create a custom resource of type `Comment`, it will publish a comment on an issue of this repository.

Go ahead, have fun! Play around with this demo and see your comment appear on Github. :smile:

## Installing

```
$ go get github.com/nikhita/kube-custom-controller/...
$ go build
```

## Usage

**Prerequisites**:

- A Github API token
- A kubeconfig file

1. Make sure you export your `kubeconfig` file and Github API token as follows:

    ```
    $ export KUBECONFIG=path/to/kubeconfig-file
    $ export TOKEN=<github-api-token>
    ```

2. Then run the controller:

    ```
    $ ./kube-custom-controller
    ```

3. Create a `CustomResourceDefinition` to register the type `Comment`.

    ```
    $ kubectl create -f artifacts/crd.yaml
    ```

4. Create a custom object of type `Comment`. Don't forget to change the message in this file!

    ```
    $ kubectl create -f artifacts/cr.yaml
    ```
