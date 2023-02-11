### k8s通过NodePort方式获取真实IP

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: conn-rpc
  name: tcp-nodeport
  namespace: xxim
spec:
  externalTrafficPolicy: Local # 保证源IP不被改变
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: 8999-8999-tcp-3ku1mxcacoy
    nodePort: 30007
    port: 8999
    protocol: TCP
    targetPort: 8999
  selector:
    app: conn-rpc
  sessionAffinity: None
  type: NodePort
```