deploy.yaml contain the pod and service resource.

Assumption: you already have a redis db running in the cluster
and you know the service info for that db. In this demo we
have that redis db running in the `demo` namespace with
service name `redis-leader` and we are deploying this
resource in the same namespace. if you have it in different
namespace then adjust the `DB_HOST` env in the pod resource
before deploying it to cluster.

```
$ oc project
Using project "demo" from context named "microshift" on server "https://api.crc.testing:6443".

$ oc get pods
NAME                              READY   STATUS    RESTARTS   AGE
redis-leader-5455744567-6w26s     1/1     Running   0          4d23h
redis-microservice                1/1     Running   0          7m8s

$ oc get svc
NAME                 TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
redis-leader         ClusterIP   10.43.51.59   <none>        6379/TCP   4d23h
redis-microservice   ClusterIP   10.43.64.73   <none>        8080/TCP   6m19s

 $ oc logs redis-microservice 
2023/04/24 07:07:22 Connected to Redis: PONG
```

Redis db which is part of this demo is deployed using https://github.com/praveenkumar/examples/blob/master/guestbook/all-in-one/guestbook-all-in-one.yaml
resource.
