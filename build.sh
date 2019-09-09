#scp -r /Users/yang/laiye/yangjinjie/kube-eventer/* root@172.16.13.103:/root/kube-eventer/kube-eventer
scp -r /Users/yang/laiye/yangjinjie/kube-eventer/sinks/* root@172.16.13.103:/root/kube-eventer/kube-eventer/sinks


# GOARCH=amd64 CGO_ENABLED=0 go build -o kube-eventer

# docker build -t liyehaha/kube-eventer -f Dockerfile  .
# docker push liyehaha/kube-eventer