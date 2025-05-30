# Online Boutique

##

This tutorial shows you how to deploy the cloud-native microservices demo application **[Online Boutique](https://github.com/go-micro/demo)** to a Kubernetes cluster.

You'll be able to run Online Boutique on:

- a free [minikube](https://minikube.sigs.k8s.io/docs/) cluster, which comes built in to the Cloud Shell instance
- a [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) cluster using a new or existing [Google Cloud Platform project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project)

Let's get started!

## Kubernetes cluster setup

Set up a Kubernetes cluster using the instructions below for either **minikube** or **GKE**.

---
### Minikube instructions

Minikube creates a local Kubernetes cluster on Cloud Shell.

1. Click <walkthrough-editor-spotlight spotlightId="minikube-status-bar">minikube</walkthrough-editor-spotlight> on the status bar located at the bottom of the editor window.

2. The command palette will prompt you to choose which minikube cluster to control. Select **minikube** and, in the next prompt, click **Start** if the cluster has not already been started. 

3. If prompted, authorize Cloud Shell to make a GCP API call with your credentials.

*It may take a few minutes for minikube to finish starting.*

Once minikube has started, you're ready to move on to the next step. 

--- 

## Run on Kubernetes

Now you can run Online Boutique on your Kubernetes cluster!

1. Launch the <walkthrough-editor-spotlight spotlightId="cloud-code-status-bar">Cloud Code menu</walkthrough-editor-spotlight> from the status bar and select <walkthrough-editor-spotlight spotlightId="cloud-code-run-on-k8s">Run on Kubernetes</walkthrough-editor-spotlight>.

2. If prompted to select a Skaffold Profile, select **[default]**.

3. Select **Yes** to confirm your current context.

4. If you're using a GKE cluster, you'll need to confirm your container image registry.

5. If prompted, authorize Cloud Shell to make a GCP API call with your credentials.

Cloud Code uses configurations defined in <walkthrough-editor-open-file filePath="skaffold.yaml">skaffold.yaml</walkthrough-editor-open-file> to build and deploy the app. *It may take a few minutes for the deploy to complete.*

6. Once the app is running, the local URLs will be displayed in the <walkthrough-editor-spotlight spotlightId="output">Output</walkthrough-editor-spotlight> terminal. 

7. To access your Online Boutique frontend service, click on the <walkthrough-spotlight-pointer spotlightId="devshell-web-preview-button" target="cloudshell">Web Preview button</walkthrough-spotlight-pointer> in the upper right of the editor window.

8. Select **Change Port** and enter '4503' as the port, then click **Change and Preview**. Your app will open in a new window. 


## Stop the app

To stop running the app: 

1. Go to the <walkthrough-editor-spotlight spotlightId="activity-bar-debug">Debug view</walkthrough-editor-spotlight> 

2. Click the 'Stop' icon. 🟥

3. Select **Yes** to clean up deployed resources. 

You can start, stop, and debug apps from the Debug view.

### Clean up

If you've deployed your app to a GKE cluster in your Google Cloud Platform project, you'll want to delete the cluster to avoid incurring charges.

1. Navigate to the <walkthrough-editor-spotlight spotlightId="activity-bar-cloud-k8s">Cloud Code - Kubernetes view</walkthrough-editor-spotlight> in the Activity bar.

2. Under the <walkthrough-editor-spotlight spotlightId="cloud-code-gke-explorer">Google Kubernetes Engine Explorer tab</walkthrough-editor-spotlight>, right-click on your cluster and select **Delete Cluster**.


## Conclusion

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

Congratulations! You've successfully deployed Online Boutique using Cloud Shell.

<walkthrough-inline-feedback></walkthrough-inline-feedback>

##### What's next?

Try other deployment options for Online Boutique:
- **Workload Identity**: <walkthrough-editor-open-file filePath="./docs/workload-identity.md">See these instructions</walkthrough-editor-open-file>.
- **Istio**: <walkthrough-editor-open-file filePath="./docs/service-mesh.md">See these instructions</walkthrough-editor-open-file>.
- **Anthos Service Mesh**: ASM requires Workload Identity to be enabled in your GKE cluster. <walkthrough-editor-open-file filePath="./docs/workload-identity.md">See these instructions</walkthrough-editor-open-file> to configure and deploy the app. Then, use the <walkthrough-editor-open-file filePath="./docs/service-mesh.md">service mesh guide</walkthrough-editor-open-file>.
- **Memorystore**: <walkthrough-editor-open-file filePath="./docs/memorystore.md">See these instructions</walkthrough-editor-open-file> to replace the in-cluster `redis` database with hosted Google Cloud Memorystore (redis).

Learn more about the [Cloud Shell](https://cloud.google.com/shell) IDE environment and the [Cloud Code](https://cloud.google.com/code) extension.
