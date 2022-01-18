# mongodb-probe

Fair warning, you probably should be looking at percona/mongodb_exporter rather than this project.

This project is a simple probe that exports prometheus metrics for mongodb. It connects to the nodes and exports the state and replication lag.

Proper MongoDB support requires a considerable amount of effort, this project is used internally, we only keep track of the version of mongo we use and we consider mongo the legacy database, which is why I do not plan on making this a proper open-source project with support.

Regardless, if you, like us, was frustated by mongodb_exporter, feel free to use it. A quick read on main.go will be enough to understand the env options.