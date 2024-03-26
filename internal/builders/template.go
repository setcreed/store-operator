package builders

// store-core的deployment模版
const deployTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: store-{{ .Name }}
  namespace: {{ .Namespace}}
spec:
  selector:
    matchLabels:
      app: store-{{ .Namespace}}-{{ .Name }}
  replicas: 1
  template:
    metadata:
      labels:
        app: store-{{ .Namespace}}-{{ .Name }}
        version: v1
      annotations:
        store.config/md5: ''
    spec:
      initContainers:
        - name: init-test
          image: busybox:1.28
          command: ['sh', '-c', 'echo sleeping && sleep 15']
      containers:
        - name: store-{{ .Namespace}}-{{ .Name }}-container
          image: docker.io/setcreed/store-core:v0.1
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: configdata
              mountPath: /etc/dbcore/config.yaml
              subPath: config.yaml
          ports:
             - containerPort: 8080
             - containerPort: 8090
      volumes:
       - name: configdata
         configMap:
          defaultMode: 0644
          name: store-{{ .Name }}
`

// store-core configmap 模版
const cmTemplate = `
  default:
    mode: debug
    app:
      rpcPort: 8080
      httpPort: 8090
  dbConfig:
    dsn: "[[ .Dsn ]]"
    maxOpenConn: [[ .MaxOpenConn ]]
    maxLifeTime: [[ .MaxLifeTime ]]
    maxIdleConn: [[ .MaxIdleConn ]]
  sqlConfig:
  - name: test
    sql: "select * from test"
`
