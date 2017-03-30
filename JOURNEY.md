[![Build Status](https://travis-ci.org/IBM/kubernetes-container-service-gitlab-sample.svg?branch=master)](https://travis-ci.org/IBM/kubernetes-container-service-gitlab-sample)

# Implementing Common Transactions on IBM Blockchain

## Overview
This project shows how to perform traditional data store transactions on IBM Blockchain. This surfaces as a web-based, to-do list application, allowing browse, read, edit, add, and delete (BREAD) operations.

IBM Blockchain is powered by the Linux Foundation Hyperledger Fabric 0.6 currently. That service is being deprecated, and will be replaced by the Hyperledger Fabric 1.0 architecture in the near future. At that point, this example will be updated to reflect migration patterns for businesses currently using the 0.6 architecture (privately or hosted on IBM Bluemix).

The to-do list application presented here is designed to help developers understand how common transactions needed by business processes can be adapted to use Blockchain. *Blockchain != Bitcoin.* It might be said that Bitcoin is the first Blockchain application. As a distributed ledger, the distinct characteristics such as decentralization, consensus, and encryption have broad-reaching implications to many business verticals including finance, transportation, health care, and many others.

![Flow](images/gitlab_container.png)

## Included Components
- IBM Blockchain
- OpenWhisk

## Prerequisites

Create an instance of IBM Blockchain on IBM Bluemix.

If you have not provisioned services on IBM Bluemix before, please follow the Setting Up IBM Blockchain tutorial.

You will also need a public GitHub repository.

## Steps

1. Deploy the compiled chaincode to a public GitHub repository
2. Deploy the supporting OpenWhisk action
3. Run a local web server, or run the included Node.js server
4. Using the to-do list application

# 1. Deploy the compiled chaincode to a public GitHub repository

A compiled chaincode application is included. By placing this file in a public GitHub repository, you make it available for IBM Blockchain consumption. Once the chaincode is deployed in a public GitHub repository, log into IBM Blockchain, and navigate to the deployment screen.

> Chaincode is also commonly referred to as a "smart contract" in Blockchain terminology. Chaincode is authored using the Go language. This represents the business logic of what transactions can take place on the Blockchain.

TODO: Deploy the chaincode and make note of the resulting chaincode ID.

# 2. Deploy the supporting OpenWhisk action

A JavaScript-based OpenWhisk action is included for the purposes of natural language processing (NLP) of task due dates. This allows the application to function without a calendar component, allowing the user to input dates in the form of "today" or "one week from today" as well as traditional "April 1".

> OpenWhisk is a serverless deployment option, also commonly referred to as "functions as a service". This deployment option allows server-side business logic implementation, without the need of running a dedicated server instance.

TODO: Deploy the OpenWhisk action using the command-line tooling.

# 3. Run a local web server, or run the included Node.js server

In order for the web-based to-do list application to work, it must be run from a web server. This server does not need to be publicly available in order for the application to function. On Mac, a common approach is to use the built-in PHP installation to run an in-place web server.

At the command line, navigate to the "/bluemix/public" directory. Launch the PHP web server in-place.

```bash
php -S localhost:8081
```

Alternatively, a basic Node.js Express web server is provided. To use this as your web server, you must first install the Node.js dependencies. Dependencies are defined in "/bluemix/package.json". Navigate to the "/bluemix" directory, install the dependencies, and then run the application.

```bash
npm install
node app.js
```

> IBM Blockchain supports cross-origin resource sharing (CORS). The result is that the browser can communicate directly to IBM Blockchain without the need for a proxy server.

# 4. Using the to-do list application




First, install [Docker CLI](https://www.docker.com/community-edition#/download).

Then, install the Bluemix container registry plugin.

```bash
bx plugin install container-registry -r bluemix
```

Once the plugin is installed you can log into the Bluemix Container Registry.

```bash
bx cr login
```

If this is the first time using the Bluemix Container Registry you must set a namespace which identifies your private Bluemix images registry. It can be between 4 and 30 characters.

```bash
bx cr namespace-add <namespace>
```

Verify that it works.

```bash
bx cr images
```


# 2. Build PostgreSQL and GitLab containers

PostgreSQL and GitLab containers need to be built. Redis container can be used as is from Docker Hub

Build the PostgreSQL container.

```bash
cd containers/postgres
docker build -t registry.ng.bluemix.net/<namespace>/gitlab-postgres .
docker push registry.ng.bluemix.net/<namespace>/gitlab-postgres
```

Build the GitLab container.

```bash
cd containers/gitlab
docker build -t registry.ng.bluemix.net/<namespace>/gitlab .
docker push registry.ng.bluemix.net/<namespace>/gitlab
```


After finish building the images in bluemix registery, please modify the container images in your yaml files. 

i.e. 
1. In postgres.yaml, change `docker.io/tomcli/postgres:latest` to `registry.ng.bluemix.net/<namespace>/gitlab-postgres`
2. In gitlab.yaml, change `docker.io/tomcli/gitlab:latest` to `registry.ng.bluemix.net/<namespace>/gitlab`

> Note: Replace `<namespace>` to your own container registry namespace. You can check your namespace via `bx cr namespaces`

# 3. Create Services and Deployments

Run the following commands or run the quickstart script `bash quickstart.sh` with your Kubernetes cluster.

```bash
kubectl create -f local-volumes.yaml
kubectl create -f postgres.yaml
kubectl create -f redis.yaml
kubectl create -f gitlab.yaml
```

After you have created all the services and deployments, wait for 3 to 5 minutes. You can check the status of your deployment on Kubernetes UI. Run 'kubectl proxy' and go to URL 'http://127.0.0.1:8001/ui' to check when the GitLab container becomes ready.

![Kubernetes Status Page](images/kube_ui.png)

After few minutes the following commands to get your public IP and NodePort number.

```bash
$ kubectl get nodes
NAME             STATUS    AGE
169.47.241.106   Ready     23h
$ kubectl get svc gitlab
NAME      CLUSTER-IP     EXTERNAL-IP   PORT(S)                     AGE
gitlab    10.10.10.148   <nodes>       80:30080/TCP,22:30022/TCP   2s
```

> Note: The 30080 port is for gitlab UI and the 30022 port is for ssh.

Congratulation. Now you can use the link **http://[IP]:30080** to access your gitlab site on browser.

> Note: For the above example, the link would be http://169.47.241.106:30080  since its IP is 169.47.241.106 and the UI port number is 30080. 


# 4. Using GitLab
Now that Gitlab is running you can register as a new user and create a project.

![Registration page](images/register.png)


After logging in as your newly-created user you can create a new project.

![Create project](images/new_project.png)

Once a project has been created you'll be asked to add an SSH key for your user.

To verify that your key is working correctly run:

```bash
ssh -T git@<IP> -p 30022
```

Which should result in:

```bash
Welcome to GitLab, <user>!
```

Now you can clone your project.
```bash
git clone ssh://git@<IP>:30022/<user>/<project name>
```

Add a file and commit:
```bash
echo "Gitlab project" > README.md
git add README.md
git commit -a -m "Initial commit"
git push origin master
```

You can now see it in the Gitlab UI.
![Repo](images/first_commit.png)

# Troubleshooting
If a pod doesn't start examine the logs.
```bash
kubectl get pods
kubectl logs <pod name>
```


To delete all your services, deployments, and persistent volume claim, run

```bash
kubectl delete deployment,service,pvc -l app=gitlab
```

To delete your persistent volume, run

```bash
kubectl delete pv local-volume-1 local-volume-2 local-volume-3
```

# License
[Apache 2.0](LICENSE.txt)