# Development Guide 

This doc explains how to build and run the OnlineBoutique source code locally using the `skaffold` command-line tool.  

## Prerequisites 

- [Docker for Desktop](https://www.docker.com/products/docker-desktop).
- kubectl (can be installed via `gcloud components install kubectl`)
- [skaffold **1.27+**](https://skaffold.dev/docs/install/) (latest version recommended), a tool that builds and deploys Docker images in bulk. 
- A Google Cloud Project with Google Container Registry enabled. 
- Enable GCP APIs for Cloud Monitoring, Tracing, Debugger, Profiler:
```
gcloud services enable monitoring.googleapis.com \
    cloudtrace.googleapis.com \
    clouddebugger.googleapis.com \
    cloudprofiler.googleapis.com
```
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (optional - see Local Cluster)
- [Kind](https://kind.sigs.k8s.io/) (optional - see Local Cluster)

## Local Cluster 

1. Launch a local Kubernetes cluster with one of the following tools:

    - To launch **Minikube** (tested with Ubuntu Linux). Please, ensure that the
       local Kubernetes cluster has at least:
        - 4 CPUs
        - 4.0 GiB memory
        - 32 GB disk space

      ```shell
      minikube start --cpus=4 --memory 4096 --disk-size 32g
      ```

    - To launch **Docker for Desktop** (tested with Mac/Windows). Go to Preferences:
        - choose “Enable Kubernetes”,
        - set CPUs to at least 3, and Memory to at least 6.0 GiB
        - on the "Disk" tab, set at least 32 GB disk space

    - To launch a **Kind** cluster:

      ```shell
      kind create cluster
      ```

2. Run `kubectl get nodes` to verify you're connected to the respective control plane.

3. Run `skaffold run` (first time will be slow, it can take ~20 minutes).
   This will build and deploy the application. If you need to rebuild the images
   automatically as you refactor the code, run `skaffold dev` command.

4. Run `kubectl get pods` to verify the Pods are ready and running.

5. Access the web frontend through your browser
    - **Minikube** requires you to run a command to access the frontend service:

    ```shell
    minikube service frontend-external
    ```

    - **Docker For Desktop** should automatically provide the frontend at http://localhost:80

    - **Kind** does not provision an IP address for the service.
      You must run a port-forwarding process to access the frontend at http://localhost:8080:

    ```shell
    kubectl port-forward deployment/frontend 8080:8080
    ```

## Cleanup

If you've deployed the application with `skaffold run` command, you can run
`skaffold delete` to clean up the deployed resources.
