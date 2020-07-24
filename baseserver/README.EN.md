# TARS Baseserver Deployment

## Instructions

1. Prerequisites 
      - Install Kubernetes; you may use kubectl or other tools of your choice to control the clusters
      - An endpoint for docker commands

2. Deploying tars db

      ```
      // acquire relevant documents on deployment under baseserver file
      cd baseserver
      make deploy

      // Create a MySQL database (in practice, a cloud db is recommended)  
      kubectl apply -f yaml/db_all_in_one.yaml

      // Acquire the names of the database with following command 
      kubectl get pods | grep tars-db

      // Edit the name in db/install_db_k8s.sh, and then import the data
      sh db/install_db_k8s.sh
      ```
      
      For existing tars db, executing sql files may erase original data. You only need to import the absent db.

3. tars registry 

   Use `kubectl apply -f yaml/registry.yaml` to deploy tars registry. 
   If your k8s was not used in the process of creating the db, please edit the data path in the `registry.yaml` file accordingly.

4. tarsweb

   Use `kubectl apply -f yaml/tarsweb.yaml` to deploy. 
   By default tarsweb uses port 3000，Or you may use a method of your choice, and check the status on your browser.
   
   Note: when the current tarsweb version is not compatible with k8s' scenarios, the page has restart/stop options, but they will not execute successfully.

5. Deploying other services

   In `tarsnotify` for instance, please use `kubectl apply -f yaml/registry.yaml` to deploy. All db configurations can be edited accordingly. 
    All other services can be deployed in this fashion，as you may simply change `tarsnotify` to your desired service. Currently available services：
      1. tarslog
      2. tarsconfig
      3. tarsproperty
      4. tarsstat
      5. tarsquerystat
      6. tarsqueryproperty
  
## Creating images 
`make registry` generates the image of registry

`make web` generates the image of tarsweb

`make img SERVER=XXX` generates the image of basic server=xxx 

Note: cppregistry is the original master control, and can be merged into registry later on. 
