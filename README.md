## simple application that can run your command and send some information about it
## the metrics are:
### push_job_status:
status code of the command 
### push_job_completionTime:
the time that the command finished 
### push_job_successTime:
the time that the command finished successfully (this metric wont create if the job was not successfull)
### push_job_duration:
duration of the task

### build:
$go mod download  
$go mod verify  
$go build -o pushjob  

### options:
-c : command to run  
-h : pushgateway host to push the metrics to  
-j : jobname (its like the promtheus jobname in its config)  
-l : labels  
-o : filename to store output log   




### example:
$./pushjob -c "rm -rf /tmp/images" -j "remove" -l "service=cdn" -l "cluster=dev" -l "frequency=daily" -l "servername=cdn-node1" -h "pushgateway.example.com" -o /var/log/log.log  
