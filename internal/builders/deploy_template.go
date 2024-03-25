package builders

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
    spec:
      containers:
        - name: store-{{ .Namespace}}-{{ .Name }}-container
          image: docker.io/setcreed/store-core:v0.1
          imagePullPolicy: IfNotPresent
          ports:
             - containerPort: 8080
             - containerPort: 8090
`
