# Local PV Provisioner BDD

Local PV Provisioner BDD tests are developed using ginkgo & gomega libraries.

## How to run the tests?

### Pre-requisites

- These tests are meant to be run in a single-node Kubernetes
  cluster with one single available blockdevice (not mounted).

- Some of the tests test features exclusive to the XFS filesystem.
  To run the XFS tests, you will require the 'xfsprogs' package
  installed on your node.

- You will require the Ginkgo binary to be able to run the tests.
  Install the latest Ginkgo binary using the following command:
  ```bash
  $ go install github.com/onsi/ginkgo/ginkgo@latest
  ```

- Get your Kubernetes Cluster ready and make sure you can run 
  kubectl from your development machine. 
  Note down the path to the `kubeconfig` file used by kubectl 
  to access your cluster.  Example: /home/\<user\>/.kube/config

- Set the KUBECONFIG environment variable on your 
  development machine to point to the kubeconfig file. 
  Example: `export KUBECONFIG=$HOME/.kube/config`

  If you do not set this ENV, you will have to pass the file 
  to the Ginkgo CLI (see below)

- The tests should not be run in parallel as it may lead to
  unavailability of blockdevices for some of the tests.

- Install required OpenEBS LocalPV Provisioner components
  Example: `kubectl apply -f https://openebs.github.io/charts/openebs-operator-lite.yaml`

### Run tests

Run the tests by being in the localpv tests folder. 
  ```bash
  $ cd <repo-directory>/tests
  $ ginkgo -v
 ```
  In case the KUBECONFIG env is not configured, you can run:
 ```bash
  $ ginkgo -v -- -kubeconfig=/path/to/kubeconfig
 ```

  If your OpenEBS LocalPV components are in a different Kubernetes namespace than 'openebs', you may use the '-openebs-namespace' flag:
 ```bash
  $ ginkgo -v -- -openebs-namespace=<your-namespace>
 ```

 > **Tip:** Raising a pull request to this repo's 'develop' branch (or any one of the release branches) will automatically run the BDD tests in GitHub Actions. You can verify your code changes by moving to the 'Checks' tab in your pull request page, and checking the results of the 'integration-test' check.